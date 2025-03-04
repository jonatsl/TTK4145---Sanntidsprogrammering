package requests

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
)

var _numFloors int = elevio.NumFloors
var _numButtons int = elevio.NumButtons

type DirnBehaviourPair struct {
	Dirn      elevio.MotorDirection
	Behaviour elevator.ElevatorBehaviour
}

func OnDoorTimeout(elev elevator.Elevator, doorTimerCh chan bool, lightsCh chan<- int, elevUpdateRealTimeCh chan<- elevator.Elevator) elevator.Elevator {

	elev = clearAtCurrentFloor(elev)
	pair := ChooseDirection(elev)

	elev.Dirn = pair.Dirn
	elev.Behaviour = pair.Behaviour
	elevUpdateRealTimeCh <- elev
	// elevator.SetAllLights(elev)
	lightsCh <- 1

	if elev.Behaviour == elevator.EB_DoorOpen {
		doorTimerCh <- true
	} else {
		elevio.SetMotorDirection(elev.Dirn)
	}

	return elev
}
// func OnRequest(elev elevator.Elevator, lightsCh chan<- int) elevator.Elevator {
// 	// switch elev.Behaviour {
// 	// case elevator.EB_DoorOpen:
// 	// 	//elev.Requests = elevator.MergeHallAndRequests(elev.Requests, HallRequests)
// 	// 	fmt.Println("on request while Door open")
// 	// case elevator.EB_Moving:
// 	// 	//elev.Requests = elevator.MergeHallAndRequests(elev.Requests, HallRequests)
// 	// 	fmt.Println("on request while moving")
// 	// case elevator.EB_Idle:
// 	//elev.Requests = elevator.MergeHallAndRequests(elev.Requests, HallRequests)
// 	if elev.Behaviour == elevator.EB_Idle {
// 		fmt.Println("on request while idle")
// 		pair := chooseDirection(elev)
// 		elev.Dirn = pair.Dirn
// 		elevio.SetMotorDirection(elev.Dirn)
// 		elev.Behaviour = pair.Behaviour
// 	}
// 	// }

// 	// elevator.SetAllLights(elev)
// 	lightsCh <- 1

// 	return elev
// }

func requests_above(e elevator.Elevator) bool {
	for f := e.Floor + 1; f < _numFloors; f++ {
		for btn := 0; btn < _numButtons; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requests_below(e elevator.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < _numButtons; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requests_here(e elevator.Elevator) bool {
	for btn := 0; btn < _numButtons; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

func ChooseDirection(e elevator.Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case elevio.MD_Up:
		if requests_above(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: elevator.EB_Moving}
		} else if requests_here(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: elevator.EB_DoorOpen}
		} else if requests_below(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: elevator.EB_Idle}
		}
	case elevio.MD_Down:
		if requests_below(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: elevator.EB_Moving}
		} else if requests_here(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: elevator.EB_DoorOpen}
		} else if requests_above(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: elevator.EB_Idle}
		}
	case elevio.MD_Stop:
		if requests_here(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: elevator.EB_DoorOpen}
		} else if requests_above(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: elevator.EB_Moving}
		} else if requests_below(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: elevator.EB_Idle}
		}
	default:
		return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: elevator.EB_Idle}
	}
}

func ShouldStop(e elevator.Elevator) bool {
	switch e.Dirn {
	case elevio.MD_Down:
		return e.Requests[e.Floor][elevio.BT_HallDown] || e.Requests[e.Floor][elevio.BT_Cab] || !requests_below(e)
	case elevio.MD_Up:
		return e.Requests[e.Floor][elevio.BT_HallUp] || e.Requests[e.Floor][elevio.BT_Cab] || !requests_above(e)
	case elevio.MD_Stop:
		return true
	default:
		return true
	}
}

// // denne skal ikke gjøre endringer på heisen, kun sjekke if should stop. kanskje gjøre heisen til const parameter?
// func ShouldClearImmediately(e elevator.Elevator, btn_floor int, btn_type elevio.ButtonType) bool { //det er vel her problemet angående at vi ikke klarer å plukke opp bestillinger i etasjen vi allerede befinner oss i?
// 	if e.Floor == btn_floor && ((e.Dirn == elevio.MD_Up && btn_type == elevio.BT_HallUp) ||
// 		(e.Dirn == elevio.MD_Down && btn_type == elevio.BT_HallDown) ||
// 		(e.Dirn == elevio.MD_Stop) ||
// 		(btn_type == elevio.BT_Cab)) {
// 		return true
// 	} else {
// 		return false
// 	} //Må man ha klammeparentes her?
// }
func ShouldClearImmediately(elev elevator.Elevator) bool {
    // if e.Floor == -1 { // Elevator uninitialized
    //     return false
    // }

    switch elev.Dirn {
    case elevio.MD_Up:
        if elev.Requests[elev.Floor][elevio.BT_HallUp] || elev.Requests[elev.Floor][elevio.BT_Cab] {
            return true
        }
    case elevio.MD_Down:
        if elev.Requests[elev.Floor][elevio.BT_HallDown] || elev.Requests[elev.Floor][elevio.BT_Cab] {
            return true
        }
    case elevio.MD_Stop:
        if elev.Requests[elev.Floor][elevio.BT_HallUp] || elev.Requests[elev.Floor][elevio.BT_HallDown] || elev.Requests[elev.Floor][elevio.BT_Cab] {
            return true
        }
    }

    return false
}

func clearAtCurrentFloor(e elevator.Elevator) elevator.Elevator {
	e.Requests[e.Floor][elevio.BT_Cab] = false
	switch e.Dirn {
	case elevio.MD_Up:
		if !requests_above(e) && !e.Requests[e.Floor][elevio.BT_HallUp] {
			e.Requests[e.Floor][elevio.BT_HallDown] = false
		}
		e.Requests[e.Floor][elevio.BT_HallUp] = false
	case elevio.MD_Down:
		if !requests_below(e) && !e.Requests[e.Floor][elevio.BT_HallDown] {
			e.Requests[e.Floor][elevio.BT_HallUp] = false
		}
		e.Requests[e.Floor][elevio.BT_HallDown] = false
	case elevio.MD_Stop:
		e.Requests[e.Floor][elevio.BT_HallUp] = false
		e.Requests[e.Floor][elevio.BT_HallDown] = false
	// default:
	// 	e.Requests[e.Floor][elevio.BT_HallUp] = false
	// 	e.Requests[e.Floor][elevio.BT_HallDown] = false
	}
	return e
}
