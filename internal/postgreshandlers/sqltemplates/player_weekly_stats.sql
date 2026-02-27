-- language=GenericSQL
INSERT INTO player_stats (player_key, game_key, week_key, stat_id, stat_value) VALUES {insert_values} ON CONFLICT DO NOTHING