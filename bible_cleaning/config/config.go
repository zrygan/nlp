package config

const (
	SRC_PATH                           = "corpus"
	CORPUS_VERSES_FOLDER               = "corpus/by_verses"
	CORPUS_SENTENCES_FOLDER            = "corpus/by_sentences"
	DST_PATH                           = "parallel_corpus"
	PARALLEL_VERSES_FOLDER             = "parallel_corpus/by_verses"
	PARALLEL_SENTENCES_FOLDER          = "parallel_corpus/by_sentences"
	WORKER_REPORT_INTERVAL_MS          = 50000 // milliseconds
	THREAD_POOL_SIZE                   = 12   // number of worker threads
	IS_DETAILED                        = false
	USE_PROGRESS_BAR                   = false
	WORKER_THREAD_REPORT_PROGRESS_RATE = 10000 // report progress every N items processed
)

const (
	TOKEN_MISSING_TRANSLATION    = "<MISSING_TRANSLATION>"
	TOKEN_NEWLINE                = "<NEWLINE>"
	TOKEN_SPACE                  = "<SPACE>"
	TOKEN_TAB                    = "<TAB>"
	TOKEN_RETURN                 = "<RETURN>"
	NGRAMS_DICE_SIMILARITY_BIAS  = 0.5
	LENGTH_RATIO_SIMILARITY_BIAS = 0.3
	PROPER_NOUNS_SIMILARITY_BIAS = 0.2
)
