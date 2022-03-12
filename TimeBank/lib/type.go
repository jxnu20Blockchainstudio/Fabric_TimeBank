package lib

//Manager 账户，虚拟管理员
type Manager struct {
	ManagerID    string  //管理员id
	ManagerAsset float64 //余额，初始化为零，用户充值的RMB总额汇总
}

//User 账户
type User struct {
	UserID             string   //用户id
	UserName           string   //用户姓名
	UserIdentification string   //用户身份证号
	Sex                string   //用户性别
	Birthday           string   //出生日期
	Address            string   //住址
	Postcode           string   //邮政编码
	Ability            []string //可提供的服务
	StarSign           int      //星级标志

	//根据个人对组织的贡献、技能、资历年限等评定
	//星级标志分为：1、一颗心；2、两颗星；3、三颗星；4、四颗星；5、五颗星
	//星越多，个人对组织的贡献越多、技能越强、资历年限越高等

	UserAsset     float64  //用户工分资产
	Comment       []string //顾客对用户服务的评价
	RecommenderID string   //推荐人ID
}

//Organization 组织，用户所属地
type Organization struct {
	OrgID      string   //组织ID，为当地邮政编码
	OrgName    string   //组织名字
	UserSum    int      //组织下拥有用户的总人数
	HaveUserID []string //组织下的所有用户的ID
}

//Servicing 进行服务
type Servicing struct {
	ServicingType    string  //服务类型
	ServicingOlderID string  //发起服务老人ID
	ServicingVolID   string  //接受服务志愿者ID
	ServicingCount   int     //服务计数
	StartTime        string  //发起服务开始时间
	ServicingValue   float64 //服务所需工分
	ServicingState   string  //服务状态
}

//ServiceTrade 交易，服务与被服务的工分交易
type ServiceTrade struct {
	TxID      string //交易ID
	ServeID   string //服务人员的ID
	CustID    string //顾客ID
	TxType    string //服务类型
	WorkHours int    //服务计数

	//如果是按时计工分，即为工作时长
	//如果是按件计工分，即为总件数
	//如果是按次计工分，即为总次数

	EndTime      string  //服务结束时间
	WorkValue    float64 //服务所需工分
	ServeComment string  //顾客评价
}

//ServiceTrade 状态
var ServiceTradingStatusConstant = func() map[string]string {
	return map[string]string{
		"require":   "请求服务", //顾客发出服务请求，等待服务人员接受
		"cancelled": "取消服务", //顾客发出服务请求后，取消服务
		//"expired":   "请求过期", //顾客发出服务请求后，无服务人员接受，服务请求过期，由所在组织安排服务人员进行服务
		"accepted":  "志愿者接受", //顾客发出请求，服务人员接受，老人等待服务上门
		"servicing": "服务中",   //服务人员正在给顾客进行服务
		"done":      "完成",    //服务人员确认接收服务工分，服务结束
	}
}

//TransferAsset 转移或继承工分
type TransferAsset struct {
	TransferID   string  //转移ID
	FromAsset    string  //工分来源ID
	ToAsset      string  //工分去向ID
	AssetValue   float64 //转移工分价值
	TransferTime string  //转移工分时间
}

//RechargeSystem 充值服务
type RechargeSystem struct {
	RechargeID    string  //充值ID
	ToUserID      string  //待充值用户的ID
	RechargeValue float64 //待充入工分总额
	RechargeTime  string  //充值时间
}

//JobPrice 服务工作评价；
//所有的养老互助行为，都用工分、积分来量化；
//工分计量参考市场价后打折30%，1分=1元钱。
type JobPrice struct {
	JobID               string  //工作ID
	JobName             string  //工作名称
	JobUnitCost         float64 //工作每一定价
	DetailedDescription string  //工作描述
}

const (
	ManagerKey        = "Manager-key"
	UserKey           = "User-key"
	OrganizationKey   = "Organization-key"
	ServicingKey      = "Servicing-key"
	ServiceTradeKey   = "ServiceTrade-key"
	TransferAssetKey  = "TransferAsset-key"
	RechargeSystemKey = "RechargeSystem-key"
	JobPriceKey       = "JobPrice-key"
)
