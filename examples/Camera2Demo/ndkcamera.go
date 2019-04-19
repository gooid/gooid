// +build android

package main

import (
	"container/list"
	"fmt"
	"log"
	"time"

	app "github.com/gooid/gooid"
	camera "github.com/gooid/gooid/camera24"
	"github.com/gooid/util"
)

func (s *CameraState) AddDevice(device *camera.Device) {
	if s.devices == nil {
		s.devices = list.New()
	}
	s.devices.PushBack(device)
}

func (s *CameraState) DelDevice(device *camera.Device) {
	util.Assert(s.devices != nil)
	elm := s.devices.Front()
	for ; elm != nil; elm = elm.Next() {
		if device == elm.Value.(*camera.Device) {
			s.devices.Remove(elm)
			return
		}
	}
}

func (s *CameraState) String() string {
	ret := ""
	for elm := s.devices.Front(); elm != nil; elm = elm.Next() {
		ret += fmt.Sprint(elm.Value, ",")
	}
	if ret != "" {
		return fmt.Sprint("{ ", s.bAvailable, ", [ ", ret, " ] }")
	} else {
		return fmt.Sprint("{ ", s.bAvailable, " }")
	}
}

type CameraManager struct {
	*camera.Manager
	cameras map[string]*CameraState
}

var gMgr *CameraManager

func CameraManagerInstance() *CameraManager {
	if gMgr == nil {
		mgr := &CameraManager{}
		mgr.Manager = camera.ManagerCreate()
		if mgr.Manager != nil {
			mgr.Manager.RegisterAvailabilityCallback(mgr)
			mgr.cameras = map[string]*CameraState{}

			gMgr = mgr
		}
	}
	return gMgr
}

func (mgr *CameraManager) Delete() {
	mgr.Manager.UnregisterAvailabilityCallback(mgr)
	mgr.Manager.Delete()
	mgr.Manager = nil
}

func (mgr *CameraManager) cameraState(id string) *CameraState {
	if ret, ok := mgr.cameras[id]; ok {
		return ret
	}
	ret := new(CameraState)
	mgr.cameras[id] = ret
	return ret
}

func (mgr *CameraManager) OnCameraAvailable(id string) {
	log.Println("... onCameraAvailable", id)
	mgr.cameraState(id).bAvailable = true
}

func (mgr *CameraManager) OnCameraUnavailable(id string) {
	log.Println("... onCameraUnavailable", id)
	mgr.cameraState(id).bAvailable = false
}

func (mgr *CameraManager) OnDisconnected(device *camera.Device) {
	log.Println("... OnDisconnected", device)
	mgr.cameraState(device.GetId()).DelDevice(device)
}

func (mgr *CameraManager) OnError(device *camera.Device, error int) {
	log.Println("... OnError", device, error)
	mgr.cameraState(device.GetId()).DelDevice(device)
}

func (mgr *CameraManager) dump() string {
	cstr := ""
	for id, s := range mgr.cameras {
		cstr += id + " = " + s.String() + ";"
	}
	if cstr == "" {
		return ""
	}
	return cstr
}

func (mgr *CameraManager) GetSensorOrientation(id string) (int, int) {
	metadata, err := mgr.GetCameraCharacteristics(id)
	util.Assert(err)
	defer metadata.Free()

	lens, err := metadata.GetConstEntry(camera.LENS_FACING)
	util.Assert(err)
	orientation, err := metadata.GetConstEntry(camera.SENSOR_ORIENTATION)
	util.Assert(err)

	return int(lens.Data().([]byte)[0]), int(orientation.Data().([]int32)[0])
}

func (mgr *CameraManager) GetSupportPixels(id string) [][2]int {
	metadata, err := mgr.GetCameraCharacteristics(id)
	util.Assert(err)
	defer metadata.Free()

	cfgs, err := metadata.GetConstEntry(camera.SCALER_AVAILABLE_STREAM_CONFIGURATIONS)
	util.Assert(err)

	// SCALER_AVAILABLE_STREAM_CONFIGURATIONS:
	//  input = entry.data.i32[i * 4 + 3];
	//  format = entry.data.i32[i * 4 + 0];
	//  width = entry.data.i32[i * 4 + 1];
	//  height = entry.data.i32[i * 4 + 2];

	ps := [][2]int{}
	ds := cfgs.Data().([]int32)
	// 去除重复的
	ms := map[string]bool{}
	for i := 0; i < cfgs.Count()/4; i++ {
		w, h := int(ds[i*4+1]), int(ds[i*4+2])
		s := fmt.Sprint(w, "x", h)
		if !ms[s] {
			ps = append(ps, [2]int{w, h})
			ms[s] = true
		}
	}
	return ps
}

func (mgr *CameraManager) OpenCamera(id string) (*NDKCamera, error) {
	device, err := mgr.Manager.OpenCamera(id, mgr)
	if err == nil {
		mgr.cameraState(id).AddDevice(device)

		return NewNDKCamera(device), nil
	}
	return nil, err
}

// Camera
const (
	PREVIEW_REQUEST_IDX = iota
	JPG_CAPTURE_REQUEST_IDX
	CAPTURE_REQUEST_COUNT
)

type SessionState int

const (
	NONE   SessionState = iota
	READY               // session is ready
	ACTIVE              // session is busy
	CLOSED              // session is closed(by itself or a new session evicts)
)

type CaptureRequestInfo struct {
	outputNativeWindow *app.Window
	sessionOutput      *camera.CaptureSessionOutput
	target             *camera.OutputTarget
	request            *camera.CaptureRequest
	template           camera.DeviceRequestTemplate
	sessionSequenceId  int
}

// NDKCamera
type NDKCamera struct {
	device *camera.Device

	outputContainer *camera.CaptureSessionOutputContainer
	captureSession  *camera.CaptureSession

	captureSessionState     SessionState
	captureSessionStateChan chan SessionState

	requests [CAPTURE_REQUEST_COUNT]CaptureRequestInfo
}

func NewNDKCamera(device *camera.Device) *NDKCamera {
	o := &NDKCamera{device: device}
	var err error
	o.outputContainer, err = camera.CaptureSessionOutputContainerCreate()
	util.Assert(err)

	o.captureSessionState = NONE
	o.captureSessionStateChan = make(chan SessionState, 1)
	return o
}

func (o *NDKCamera) Close() {
	if o == nil || o.device == nil {
		return
	}
	if o.SessionState() == ACTIVE {
		o.Repeat(&o.requests[PREVIEW_REQUEST_IDX], false)
	}

	for i, _ := range o.requests {
		o.Release(&o.requests[i])
	}

	if o.captureSession != nil {
		o.captureSession.Close()
		o.captureSession = nil
	}

	o.outputContainer.Free()
	o.outputContainer = nil

	o.device.Close()
	o.device = nil

	if o.captureSessionStateChan != nil {
		close(o.captureSessionStateChan)
		o.captureSessionStateChan = nil
	}
}

func (o *NDKCamera) addSessionOutput(out *camera.CaptureSessionOutput) error {
	if o.captureSession != nil {
		o.captureSession.Close()
		o.captureSession = nil
	}
	o.outputContainer.Add(out)

	var err error
	o.captureSession, err = o.device.CreateCaptureSession(o.outputContainer, o)
	if err == nil {
		o.captureSessionState = READY
	}
	return err
}

func (o *NDKCamera) removeSessionOutput(out *camera.CaptureSessionOutput) error {
	if o.captureSession != nil {
		o.captureSession.Close()
		o.captureSession = nil
	}
	o.outputContainer.Remove(out)

	var err error
	o.captureSession, err = o.device.CreateCaptureSession(o.outputContainer, o)
	if err == nil {
		o.captureSessionState = READY
	}
	return err
}

func (o *NDKCamera) CreateRequest(req *CaptureRequestInfo, win *app.Window, imageRotation int, template camera.DeviceRequestTemplate) *CaptureRequestInfo {
	req.outputNativeWindow = win
	req.template = template

	var err error

	req.outputNativeWindow.Acquire()
	req.sessionOutput, err = camera.CaptureSessionOutputCreate(req.outputNativeWindow)
	util.Assert(err)

	req.target, err = camera.CameraOutputTargetCreate(req.outputNativeWindow)
	util.Assert(err)
	req.request, err = o.device.CreateCaptureRequest(req.template)
	util.Assert(err)
	err = req.request.AddTarget(req.target)
	util.Assert(err)

	err = o.addSessionOutput(req.sessionOutput)
	util.Assert(err)
	return req
}

func (o *NDKCamera) Release(req *CaptureRequestInfo) {
	if req.sessionOutput != nil {
		log.Println("CaptureRequestInfo.Release:", req.template)
		o.removeSessionOutput(req.sessionOutput)

		req.request.RemoveTarget(req.target)
		req.request.Free()
		req.request = nil
		req.target.Free()
		req.target = nil

		req.sessionOutput.Free()
		req.sessionOutput = nil
		req.outputNativeWindow.Release()
		req.outputNativeWindow = nil
	}
}

func (o *NDKCamera) Repeat(req *CaptureRequestInfo, start bool) {
	if start {
		err := o.captureSession.SetRepeatingRequest([]*camera.CaptureRequest{req.request})
		if err == nil {
			o.WaitSessionState(ACTIVE)
		}
	} else if o.SessionState() == ACTIVE {
		err := o.captureSession.StopRepeating()
		if err == nil {
			o.WaitSessionState(READY)
		}
	} else {
		util.Assert(false, "captureSessionState ", o.SessionState())
	}
}

// for TakePhoto
/**
 * Process JPG capture SessionCaptureCallback_OnFailed event
 * If this is current JPG capture session, simply resume preview
 * @param session the capture session that failed
 * @param request the capture request that failed
 * @param failure for additional fail info.
 */
func (o *NDKCamera) OnCaptureFailed(_ *camera.CaptureSession,
	request *camera.CaptureRequest, failure *camera.CaptureFailure) {
	log.Println("OnCaptureFailed:", failure.Reason())
	if o != nil && request == o.requests[JPG_CAPTURE_REQUEST_IDX].request {
		util.Assert(failure.SequenceId() ==
			o.requests[JPG_CAPTURE_REQUEST_IDX].sessionSequenceId,
			"Error jpg sequence id")
		//o.StartPreview(true)
	}
}

/**
 * Process event from JPEG capture
 *    SessionCaptureCallback_OnSequenceEnd()
 *
 * If this is jpg capture, turn back on preview after a catpure.
 */
func (o *NDKCamera) OnCaptureSequenceCompleted(_ *camera.CaptureSession,
	sequenceId int, _ int64) {
	log.Println("OnCaptureSequenceCompleted:", sequenceId)
	if sequenceId != o.requests[JPG_CAPTURE_REQUEST_IDX].sessionSequenceId {
		return
	}

	//o.StartPreview(true)
}

/**
 * Process event from JPEG capture
 *    SessionCaptureCallback_OnSequenceAborted()
 *
 * If this is jpg capture, turn back on preview after a catpure.
 */
func (o *NDKCamera) OnCaptureSequenceAborted(_ *camera.CaptureSession,
	sequenceId int) {
	log.Println("OnCaptureSequenceAborted:", sequenceId)
	if sequenceId != o.requests[JPG_CAPTURE_REQUEST_IDX].sessionSequenceId {
		return
	}

	//o.StartPreview(true)
}

// TakePhoto
func (o *NDKCamera) TakePhoto() bool {
	log.Println("TakePhoto...")
	o.StartPreview(false)

	var err error
	req := &o.requests[JPG_CAPTURE_REQUEST_IDX]

	req.sessionSequenceId, err =
		o.captureSession.Capture(o, []*camera.CaptureRequest{req.request})
	util.Assert(err)
	o.WaitSessionStateTimeout(ACTIVE, time.Second)
	o.WaitSessionStateTimeout(READY, time.Second)

	o.StartPreview(true)
	return true
}

func (o *NDKCamera) resetTakePhoto(win *app.Window, imageRotation int) {
	util.Assert(o.requests[JPG_CAPTURE_REQUEST_IDX].sessionOutput == nil)
	o.CreateRequest(&o.requests[JPG_CAPTURE_REQUEST_IDX], win, imageRotation, camera.TEMPLATE_STILL_CAPTURE)

	o.requests[JPG_CAPTURE_REQUEST_IDX].request.SetEntryI32(camera.JPEG_ORIENTATION, []int32{int32(imageRotation)})
}

// preview
func (o *NDKCamera) StartPreview(start bool) {
	req := &o.requests[PREVIEW_REQUEST_IDX]
	if start && o.SessionState() == ACTIVE {
		return
	}
	if !start && o.SessionState() != ACTIVE {
		return
	}
	if start {
		util.Assert(o.SessionState() == READY)
		o.Repeat(req, start)
	} else {
		util.Assert(o.SessionState() == ACTIVE)
		o.Repeat(req, start)
	}
}

func (o *NDKCamera) releasetPreview() {
	o.Release(&o.requests[PREVIEW_REQUEST_IDX])
}

func (o *NDKCamera) resetPreview(win *app.Window, imageRotation int) {
	util.Assert(o.requests[PREVIEW_REQUEST_IDX].sessionOutput == nil)
	o.CreateRequest(&o.requests[PREVIEW_REQUEST_IDX], win, imageRotation, camera.TEMPLATE_PREVIEW)

	// Only preview request is in manual mode, JPG is always in Auto mode
	// JPG capture mode could also be switch into manual mode and control
	// the capture parameters, this sample leaves JPG capture to be auto mode
	// (auto control has better effect than author's manual control)

	o.requests[PREVIEW_REQUEST_IDX].request.SetEntryU8(camera.CONTROL_AE_MODE, []byte{camera.CONTROL_AE_MODE_OFF})
}

// Session Callbacks
func (o *NDKCamera) onSessionState(session *camera.CaptureSession, state SessionState) {
	log.Println(" .. onSessionState:", state)
	if session == nil || session != o.captureSession {
		log.Println("CaptureSession is %s", session, "NOT our session", o.captureSession)
		return
	}

	o.captureSessionState = state
	o.captureSessionStateChan <- state
}

func (o *NDKCamera) SessionState() SessionState {
	return o.captureSessionState
}

func (o *NDKCamera) WaitSessionStateTimeout(s SessionState, timeout time.Duration) {
	if timeout == 0 {
		timeout = time.Millisecond * 200
	}
	select {
	case v := <-o.captureSessionStateChan:
		log.Println("WaitSessionState:", v)
	case <-time.After(timeout):
		log.Println("WaitSessionState: timeout.")
	}
}

func (o *NDKCamera) WaitSessionState(s SessionState) {
	o.WaitSessionStateTimeout(s, 0)
}

// Session Callbacks interface
func (o *NDKCamera) OnClosed(session *camera.CaptureSession) {
	o.onSessionState(session, CLOSED)
	close(o.captureSessionStateChan)
	o.captureSessionStateChan = nil
}

func (o *NDKCamera) OnActive(session *camera.CaptureSession) {
	o.onSessionState(session, ACTIVE)
}

func (o *NDKCamera) OnReady(session *camera.CaptureSession) {
	o.onSessionState(session, READY)
}
