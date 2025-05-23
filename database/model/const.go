package model

const (
	//FolderModel             = "model"
	FolderImg = "img"

	// FileName 上传的文件名
	FileName = "file_name"
	// FileSize 上传文件的总大小
	FileSize = "file_size"
	// FileUrl 文件url
	FileUrl = "file_url"
	// FilePath 文件存储的路径
	FilePath = "file_path"
	// FileUploadSize 文件已经上传的大小
	FileUploadSize = "upload_size"
)

const (
	FileDateFormatLayout = "2006-01-02"
)

const (
	TableNameUser          = "user"
	TableNameTrain         = "train"
	TableNameStation       = "train_station"
	TableNameTrainSchedule = "train_schedule"
	TableNameTrainStop     = "train_stop"
	TableNameTrainSeat     = "train_seat"
	TableNameUserOrder     = "user_order"
	TableNameRole          = "role"
	TableNameRouters       = "routers"
	TableNameRoleRouters   = "routers_roles"
	TableNameUserRoles     = "user_roles"
	TableNameCasBinRule    = "casbin_rule"

	FieldUUID          = "uuid"
	FieldID            = "id"
	FieldName          = "name"
	FieldUserName      = "username"
	FieldUserID        = "user_id"
	FieldNickName      = "nick_name"
	FieldUserPhone     = "phone"
	FieldUserEmail     = "email"
	FieldPositionKind  = "kind"
	FieldTrainID       = "train_id"
	FieldTrainType     = "train_type"
	FieldTrainName     = "train_name"
	FieldScheduleID    = "schedule_id"
	FieldDepartureTime = "departure_date"
	FieldSeatType      = "seat_type"
	FieldSeatNums      = "seat_nums"
	FieldSeatNowNums   = "seat_now_nums"
	FieldUri           = "uri"
	FieldMethod        = "Method"
	FieldRoutersID     = "routers_id"
	FieldRoleID        = "role_id"

	FieldStartPosition   = "start"
	FieldEndPosition     = "end"
	FieldStationCode     = "code"
	FieldStationCity     = "city"
	FieldStationProvince = "province"

	FieldOrderIsPay    = "is_pay"
	FieldOrderIsDepart = "is_depart"
	FieldOrderIsCancel = "is_cancel"

	FieldCreatedTime = "created_at"
	FieldUpdatedTime = "updated_at"
	FieldDeletedTime = "deleted_at"
)

const (
	PreloadUser    = "Users"
	PreloadRoles   = "Roles"
	PreloadRouters = "Routers"
	PreloadStops   = "Stops"
	PreloadSeats   = "Seats"
)
