package roomLogic

type Room struct {
	Id string
	PeopleNum int//房间人数
	PeoplesAc []string//房间用户id
}

func NewRoom() *Room{
	return &Room{
		Id:        "",
		PeopleNum: 0,
		PeoplesAc: nil,
	}
}

/**
进入房间
 */
func (r *Room) Push(account string){
	j := false
	for _,n := range r.PeoplesAc {
		if n == account {
			j = true
		}
	}
	if !j {
		r.PeoplesAc = append(r.PeoplesAc, account)
	}
}

/**
退出房间
*/
func (r *Room) Pop(account string){
	count := len(r.PeoplesAc)
	if count == 0 {
		return
	}
	if count == 1 && r.PeoplesAc[0] == account {
		r.PeoplesAc = []string{}
	}
	for i := range r.PeoplesAc {
		if r.PeoplesAc[i] == account && i == count {
			r.PeoplesAc =  r.PeoplesAc[:count]
		}else if r.PeoplesAc[i] == account{
			r.PeoplesAc = append(r.PeoplesAc[:i],r.PeoplesAc[i+1:]...)
			break
		}
	}
}


/**
获取当前房间的所有用户account
 */
func (r *Room) GetPeople() []string {
	return r.PeoplesAc
}
