package everything

import (
	"Daemon-test/electron/utils"
	"syscall"
)

const (
	EVERYTHING_OK                                          = 0 // no error detected
	EVERYTHING_ERROR_MEMORY                                = 1 // out of memory.
	EVERYTHING_ERROR_IPC                                   = 2 // Everything search client is not running
	EVERYTHING_ERROR_REGISTERCLASSEX                       = 3 // unable to register window class.
	EVERYTHING_ERROR_CREATEWINDOW                          = 4 // unable to create listening window
	EVERYTHING_ERROR_CREATETHREAD                          = 5 // unable to create listening thread
	EVERYTHING_ERROR_INVALIDINDEX                          = 6 // invalid index
	EVERYTHING_ERROR_INVALIDCALL                           = 7 // invalid call
	EVERYTHING_ERROR_INVALIDREQUEST                        = 8 // invalid request data, request data first.
	EVERYTHING_ERROR_INVALIDPARAMETER                      = 9 // bad parameter.
	EVERYTHING_REQUEST_FILE_NAME                           = 0x00000001
	EVERYTHING_REQUEST_PATH                                = 0x00000002
	EVERYTHING_REQUEST_FULL_PATH_AND_FILE_NAME             = 0x00000004
	EVERYTHING_REQUEST_EXTENSION                           = 0x00000008
	EVERYTHING_REQUEST_SIZE                                = 0x00000010
	EVERYTHING_REQUEST_DATE_CREATED                        = 0x00000020
	EVERYTHING_REQUEST_DATE_MODIFIED                       = 0x00000040
	EVERYTHING_REQUEST_DATE_ACCESSED                       = 0x00000080
	EVERYTHING_REQUEST_ATTRIBUTES                          = 0x00000100
	EVERYTHING_REQUEST_FILE_LIST_FILE_NAME                 = 0x00000200
	EVERYTHING_REQUEST_RUN_COUNT                           = 0x00000400
	EVERYTHING_REQUEST_DATE_RUN                            = 0x00000800
	EVERYTHING_REQUEST_DATE_RECENTLY_CHANGED               = 0x00001000
	EVERYTHING_REQUEST_HIGHLIGHTED_FILE_NAME               = 0x00002000
	EVERYTHING_REQUEST_HIGHLIGHTED_PATH                    = 0x00004000
	EVERYTHING_REQUEST_HIGHLIGHTED_FULL_PATH_AND_FILE_NAME = 0x00008000
	EVERYTHING_SORT_NAME_ASCENDING                         = 1
	EVERYTHING_SORT_NAME_DESCENDING                        = 2
	EVERYTHING_SORT_PATH_ASCENDING                         = 3
	EVERYTHING_SORT_PATH_DESCENDING                        = 4
	EVERYTHING_SORT_SIZE_ASCENDING                         = 5
	EVERYTHING_SORT_SIZE_DESCENDING                        = 6
	EVERYTHING_SORT_EXTENSION_ASCENDING                    = 7
	EVERYTHING_SORT_EXTENSION_DESCENDING                   = 8
	EVERYTHING_SORT_TYPE_NAME_ASCENDING                    = 9
	EVERYTHING_SORT_TYPE_NAME_DESCENDING                   = 10
	EVERYTHING_SORT_DATE_CREATED_ASCENDING                 = 11
	EVERYTHING_SORT_DATE_CREATED_DESCENDING                = 12
	EVERYTHING_SORT_DATE_MODIFIED_ASCENDING                = 13
	EVERYTHING_SORT_DATE_MODIFIED_DESCENDING               = 14
	EVERYTHING_SORT_ATTRIBUTES_ASCENDING                   = 15
	EVERYTHING_SORT_ATTRIBUTES_DESCENDING                  = 16
	EVERYTHING_SORT_FILE_LIST_FILENAME_ASCENDING           = 17
	EVERYTHING_SORT_FILE_LIST_FILENAME_DESCENDING          = 18
	EVERYTHING_SORT_RUN_COUNT_ASCENDING                    = 19
	EVERYTHING_SORT_RUN_COUNT_DESCENDING                   = 20
	EVERYTHING_SORT_DATE_RECENTLY_CHANGED_ASCENDING        = 21
	EVERYTHING_SORT_DATE_RECENTLY_CHANGED_DESCENDING       = 22
	EVERYTHING_SORT_DATE_ACCESSED_ASCENDING                = 23
	EVERYTHING_SORT_DATE_ACCESSED_DESCENDING               = 24
	EVERYTHING_SORT_DATE_RUN_ASCENDING                     = 25
	EVERYTHING_SORT_DATE_RUN_DESCENDING                    = 26
)

var everything = syscall.NewLazyDLL("everything/dll/Everything64.dll")

func SetSearch(text string) {
	setSearch := everything.NewProc("Everything_SetSearchW")
	setSearch.Call(utils.Str2Ptr(text))
}

func Query(isWait bool) bool {
	query := everything.NewProc("Everything_QueryW")
	ret, _, _ := query.Call(utils.Bool2Ptr(isWait))
	return utils.Ptr2Bool(ret)
}

func SetRequestFlags(flag int) {
	setRequestFlags := everything.NewProc("Everything_SetRequestFlags")
	setRequestFlags.Call(utils.Int2Ptr(flag))
}

func SetSort(flag int) {
	setSort := everything.NewProc("Everything_SetSort")
	setSort.Call(utils.Int2Ptr(flag))
}

func GetNumResults() int {
	getNumResults := everything.NewProc("Everything_GetNumResults")
	ret, _, _ := getNumResults.Call()
	return utils.Ptr2Int(ret)
}

func GetResultFileName(index int) string {
	getResultFileName := everything.NewProc("Everything_GetResultFileNameW")
	ret, _, _ := getResultFileName.Call(utils.Int2Ptr(index))
	return utils.Ptr2Str(ret)
}

func GetResultPathW(index int) string {
	getResultPathW := everything.NewProc("Everything_GetResultPathW")
	ret, _, _ := getResultPathW.Call(utils.Int2Ptr(index))
	return utils.Ptr2Str(ret)
}

func SetRegex(isEnable bool) {
	setRegex := everything.NewProc("Everything_SetRegex")
	setRegex.Call(utils.Bool2Ptr(isEnable))
}
