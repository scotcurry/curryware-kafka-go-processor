package kafkahandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	"curryware-kafka-go-processor/internal/fantasyclasses/playerclasses"
	"curryware-kafka-go-processor/internal/fantasyclasses/statsclasses"
	"curryware-kafka-go-processor/internal/fantasyclasses/teamclasses"
	"curryware-kafka-go-processor/internal/fantasyclasses/transactionclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"curryware-kafka-go-processor/internal/logging"
	"curryware-kafka-go-processor/internal/postgreshandlers"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	ddkafka "github.com/DataDog/dd-trace-go/contrib/confluentinc/confluent-kafka-go/kafka.v2/v2"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	// "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	// ddkafka "gopkg.in/DataDog/dd-trace-go.v1/contrib/confluentinc/confluent-kafka-go/kafka.v2"
)

func ConsumeMessages(server string) {

	// Logging code for Datadog.
	_, err := ValidateDNSResolution(server, "")
	logging.LogInfo("Launching ConsumeMessages - curryware-kafka-go-processor on server: " + server)

	// Log the topics that exist at startup. This is informational only — the regex subscription
	// below also picks up topics created while the service is running.
	currentTopics, err := GetTopicNames(server)
	if err != nil {
		logging.LogError("Error getting topic names from Kafka server", "error", err.Error())
	}
	for i := 0; i < len(currentTopics); i++ {
		fmt.Println(fmt.Sprintf("Consuming Message(s): %s", currentTopics[i]))
	}

	// Builds the consumer. Group ID will change for different types of statistics.
	consumer, err := ddkafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  server,
		"group.id":           "curryware-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "true",
		// How often the client refreshes topic metadata. Newly created topics that match the
		// regex subscription are picked up on the next refresh (default is 5 minutes).
		"topic.metadata.refresh.interval.ms": 60000,
	}, ddkafka.WithDataStreams())
	if err != nil {
		logging.LogError("Error building consumer")
		fmt.Println("Error building consumer")
		panic(err)
	}
	defer func(consumer *ddkafka.Consumer) {
		closeErr := consumer.Close()
		if closeErr != nil {
			logging.LogError("Error closing consumer")
		}
	}(consumer)

	// Subscribing with a regex (a topic starting with "^") makes the client re-match the pattern
	// on every metadata refresh, so topics created after startup are consumed automatically.
	// Topics starting with "_" (internal topics) are excluded, matching GetTopicNames.
	err = consumer.SubscribeTopics([]string{"^[^_].*"}, nil)
	if err != nil {
		logging.LogError("Error subscribing to topics", "error", err.Error())
		panic(err)
	}
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// This is where all topic handlers get built out.  This makes it easier to add them as they come online.
	// see the bottom of this file for implementations.
	topicHandlers := map[string]func(*kafka.Message){
		"AllAvailablePlayersTopic":   processAllAvailablePlayersTopic,
		"AllLeagueInformationTopic":  processAllLeagueInformationTopic,
		"AllTeamInformationTopic":    processAllLeagueTeamInformationTopic,
		"DatadogValidationTopic":     processDatadogValidationTopic,
		"PlayerTopicDaily":           processPlayerTopicDaily,
		"StatDescriptionTopic":       processStatDescriptionTopic,
		"StatisticsTopic":            processStatisticsTopic,
		"TeamRosterInformationTopic": processTeamRosterInformationTopic,
		"TransactionTopic":           processTransactionTopic,
	}

	// This is the loop that will run forever.  Need to use Datadog to see how much processor this actually takes.
	logging.LogInfo("Jumping into the event loop")

	// This is just a timer that will post a log entry, so I know the services is running.
	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()

	run := true
	for run {
		select {
		case <-ticker.C:
			logging.LogInfo("Processing event consumer loop")
		case sig := <-signalChannel:
			fmt.Printf("Caught signal %v, exiting\n", sig)
			run = false
		default:
			event, eventError := consumer.ReadMessage(20 * time.Second)
			logging.LogDebug("Executing Event Loop")
			if eventError != nil {
				continue
			} else {
				logging.LogInfo("Event received", "bytes", len(event.Value))
			}

			topic := *event.TopicPartition.Topic
			if handler, ok := topicHandlers[topic]; ok {
				logging.LogInfo(fmt.Sprintf("Processing topic %s", topic))
				handler(event)
			} else {
				logging.LogError(fmt.Sprintf("No handler for topic %s", topic))
			}
		}
	}
}

func processPlayerTopicDaily(event *kafka.Message) {

	logging.LogInfo("Processing PlayerTopicDaily")
	playerPackage := string(event.Value)
	logging.LogInfo("Player package length: ", len(playerPackage))
	players, err := jsonhandlers.ParseJSON[[]playerclasses.PlayerInfo](playerPackage)
	if err != nil {
		logging.LogError("Error parsing player info")
		return
	}
	postgreshandlers.InsertPlayerRecord(players)
}

func processAllAvailablePlayersTopic(event *kafka.Message) {

	logging.LogInfo("Processing AllAvailablePlayersTopic")
	availablePlayersPackage := string(event.Value)
	availablePlayers, err := jsonhandlers.ParseJSON[[]playerclasses.PlayerInfo](availablePlayersPackage)
	if err != nil {
		logging.LogError("Error parsing available player info")
		return
	}

	availablePlayerCount := postgreshandlers.InsertAvailablePlayerRecord(availablePlayers)
	logging.LogInfo("Available player records inserted", "count", availablePlayerCount)
	logging.LogInfo("Available players package length", "length", len(availablePlayersPackage))
}

func processDatadogValidationTopic(event *kafka.Message) {

	logging.LogInfo("Processing DatadogValidationTopic")
	dataValidationPackage := string(event.Value)
	logging.LogInfo("Data validation package length: ", len(dataValidationPackage))
}

func processTransactionTopic(event *kafka.Message) {

	logging.LogInfo("Processing TransactionTopic")
	transactionPackage := string(event.Value)
	transactionJson, err := jsonhandlers.ParseJSON[transactionclasses.TransactionInfoWithCount](transactionPackage)
	if err != nil {
		logging.LogError("Error parsing transaction info")
		return
	}
	transactionCount := postgreshandlers.ProcessTransactionInfo(transactionJson)
	logging.LogDebug("NEEDS TO BE UPDATED: Transaction count: ", transactionCount)
	logging.LogInfo("Transaction package length: ", len(transactionPackage))
}

func processStatisticsTopic(event *kafka.Message) {

	logging.LogInfo("Processing StatisticsTopic")
	statsPackage := string(event.Value)
	statisticsInfo, err := jsonhandlers.ParseJSON[[]statsclasses.PlayerWeeklyStatsInfo](statsPackage)
	if err != nil {
		logging.LogError("Error parsing stats info")
		return
	}

	statsCount := postgreshandlers.InsertPlayerWeeklyStats(statisticsInfo)
	logging.LogInfo("Player stats inserted", "count", statsCount)
	logging.LogInfo("Statistics package length", "length", len(statsPackage))
}

func processStatDescriptionTopic(event *kafka.Message) {

	logging.LogInfo("Processing StatDescriptionTopic")
	statDescriptionPackage := string(event.Value)
	statDescriptions, err := jsonhandlers.ParseJSON[[]leagueclasses.LeagueStatDescriptionInfo](statDescriptionPackage)
	if err != nil {
		logging.LogError("Error parsing stat description info")
		return
	}

	statDescriptionCount := postgreshandlers.InsertLeagueStatInfo(statDescriptions)
	logging.LogInfo("League stat descriptions inserted", "count", statDescriptionCount)
	logging.LogInfo("Stat description package length", "length", len(statDescriptionPackage))
}

func processAllLeagueInformationTopic(event *kafka.Message) {

	logging.LogInfo("Processing AllLeagueInformationTopic")
	leagueInfoPackage := string(event.Value)
	leagueInfo, err := jsonhandlers.ParseJSON[[]leagueclasses.LeagueInformation](leagueInfoPackage)
	if err != nil {
		logging.LogError("Error parsing league information")
		return
	}

	leagueInfoCount := postgreshandlers.InsertLeagueInformation(leagueInfo)
	logging.LogInfo("League information records inserted", "count", leagueInfoCount)
	logging.LogInfo("League information package length", "length", len(leagueInfoPackage))
}

func processTeamRosterInformationTopic(event *kafka.Message) {

	logging.LogInfo("Processing TeamRosterInformationTopic")
	teamRosterPackage := string(event.Value)
	teamRosterInfo, err := jsonhandlers.ParseJSON[[]teamclasses.TeamRosterInfo](teamRosterPackage)
	if err != nil {
		logging.LogError("Error parsing team roster info")
		return
	}

	teamRosterCount := postgreshandlers.InsertTeamRosterInfo(teamRosterInfo)
	logging.LogInfo("Team roster records inserted", "count", teamRosterCount)
	logging.LogInfo("Team roster package length", "length", len(teamRosterPackage))
}

func processAllLeagueTeamInformationTopic(event *kafka.Message) {

	logging.LogInfo("Processing AllLeagueTeamInformationTopic")
	allTeamInfoPackage := string(event.Value)
	allTeamInfo, err := jsonhandlers.ParseJSON[[]leagueclasses.TeamInformation](allTeamInfoPackage)
	if err != nil {
		logging.LogError("Error parsing all team information")
		return
	}

	allTeamInfoCount := postgreshandlers.InsertTeamInformation(allTeamInfo)
	logging.LogInfo("All team information records inserted", "count", allTeamInfoCount)
	logging.LogInfo("All team information package length", "length", len(allTeamInfoPackage))
}
