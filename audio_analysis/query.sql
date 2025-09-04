-- DuckDB Query
-- FOR: stapesai/ssi-speech-emotion-recognition (emotional speech database)
-- Available at: huggingface.co/datasets/stapesai/ssi-speech-emotion-recognition

SELECT * 
  FROM train 
  WHERE emotion = 
    -- 'HAP' -- Happy
    -- 'SAD' -- Sad
    -- 'NEU' -- Neutral
    -- 'FEA' -- Fear
    -- 'ANG' -- Anger
  AND text = 'That Is Exactly what happened'
  ORDER BY RANDOM()
  LIMIT 3;