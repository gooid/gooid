package app

import (
	"fmt"
	"log"
)

//
func (ctx *Context) String() string {
	return fmt.Sprintf("Context#%s{Destory: %v, Focus: %v}",
		ctx.className, ctx.willDestory, ctx.isFocus)
}

func (act *Activity) String() string {
	return fmt.Sprintf("Activity#%p", act)
}

func (win *Window) String() string {
	return fmt.Sprintf("Window#%p", win)
}

func (q *InputQueue) String() string {
	return fmt.Sprintf("InputQueue#%p", q)
}

func (key *KeyEvent) String() string {
	var action string
	switch key.GetAction() {
	case KEY_EVENT_ACTION_DOWN:
		action = "DOWN"
	case KEY_EVENT_ACTION_UP:
		action = "UP"
	case KEY_EVENT_ACTION_MULTIPLE:
		action = "MULTIPLE"
	default:
		action = fmt.Sprint(key.GetAction())
	}
	return fmt.Sprintf("KEY, %d, %s", key.GetKeyCode(), action)
}

func (mot *MotionEvent) String() string {
	var action string
	switch mot.GetAction() & MOTION_EVENT_ACTION_MASK {
	case MOTION_EVENT_ACTION_DOWN:
		action = "DOWN"
	case MOTION_EVENT_ACTION_UP:
		action = "UP"
	case MOTION_EVENT_ACTION_MOVE:
		action = "MOVE"
	case MOTION_EVENT_ACTION_POINTER_DOWN:
		action = "POINTER_DOWN"
	case MOTION_EVENT_ACTION_POINTER_UP:
		action = "POINTER_UP"
	default:
		action = fmt.Sprint(mot.GetAction())
	}
	index := (mot.GetAction() & MOTION_EVENT_ACTION_POINTER_INDEX_MASK) >> MOTION_EVENT_ACTION_POINTER_INDEX_SHIFT
	return fmt.Sprintf("MOTION (%v), %v, %v, %v",
		index, action, int(mot.GetX(index)), int(mot.GetY(index)))
}

func (e *InputEvent) String() string {
	str := "InputEvent#{ "
	switch e.GetType() {
	case INPUT_EVENT_TYPE_KEY:
		str += (*KeyEvent)(e).String() + " }"

	case INPUT_EVENT_TYPE_MOTION:
		str += (*MotionEvent)(e).String() + " }"

	default:
		str += fmt.Sprintf("UNKONW, %p }", e)
	}

	return str
}

func (e *InputEvent) Key() *KeyEvent {
	if e.GetType() == INPUT_EVENT_TYPE_KEY {
		return (*KeyEvent)(e)
	}
	return nil
}

func (e *InputEvent) Motion() *MotionEvent {
	if e.GetType() == INPUT_EVENT_TYPE_MOTION {
		return (*MotionEvent)(e)
	}
	return nil
}

func (act *Activity) Context() *Context {
	return (*Context)(act.Instance())
}

// assert
func assert(con interface{}, infos ...interface{}) {
	if con == nil {
		return
	}
	switch obj := con.(type) {
	case bool:
		if !obj {
			info(append([]interface{}{"ASSERT:"}, infos...)...)
		}

	case error:
		info(append([]interface{}{"ASSERT:(", con, "):"}, infos...)...)

	default:
		info("ASSERT:(", con, "):", "condition must is bool/error.")
	}
}

// fatal
func fatal(v ...interface{}) {
	log.Fatal(v...)
}

// info
func info(v ...interface{}) {
	log.Println(v...)
}
