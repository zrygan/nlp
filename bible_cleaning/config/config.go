package config

const (
	SRC_PATH                           = "corpus"
	CORPUS_VERSES_FOLDER               = "corpus/by_verses"
	CORPUS_SENTENCES_FOLDER            = "corpus/by_sentences"
	DST_PATH                           = "parallel_corpus"
	PARALLEL_VERSES_FOLDER             = "parallel_corpus/by_verses"
	PARALLEL_SENTENCES_FOLDER          = "parallel_corpus/by_sentences"
	WORKER_REPORT_INTERVAL_MS          = 500 // milliseconds
	THREAD_POOL_SIZE                   = 8   // number of worker threads
<<<<<<< HEAD
	IS_DETAILED                        = true
=======
	IS_DETAILED                        = false
>>>>>>> 7ecbd5fd670f8dfc71a1b8d24ef2f53232acf56c
	USE_PROGRESS_BAR                   = true
	WORKER_THREAD_REPORT_PROGRESS_RATE = 50 // report progress every N items processed
)

const (
<<<<<<< HEAD
	TOKEN_MISSING_TRANSLATION = "<MISSING_TRANSLATION>"
	TOKEN_NEWLINE             = "<NEWLINE>"
	TOKEN_SPACE               = "<SPACE>"
	TOKEN_TAB                 = "<TAB>"
	TOKEN_RETURN              = "<RETURN>"
	DICE_SIMILARITY_THRESHOLD = 0.6
	LENGTH_RATIO_BIAS         = 1.0 - DICE_SIMILARITY_THRESHOLD
=======
	TOKEN_MISSING_TRANSLATION    = "<MISSING_TRANSLATION>"
	TOKEN_NEWLINE                = "<NEWLINE>"
	TOKEN_SPACE                  = "<SPACE>"
	TOKEN_TAB                    = "<TAB>"
	TOKEN_RETURN                 = "<RETURN>"
	NGRAMS_DICE_SIMILARITY_BIAS  = 0.5
	LENGTH_RATIO_SIMILARITY_BIAS = 0.3
	PROPER_NOUNS_SIMILARITY_BIAS = 0.2
>>>>>>> 7ecbd5fd670f8dfc71a1b8d24ef2f53232acf56c
)
