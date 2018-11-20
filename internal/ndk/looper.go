package app

/*
#include <android/looper.h>
extern int cgoLooperCallback(int fd, int events, void* data);
*/
import "C"

import "unsafe"

const (
	/**
	 * Option for ALooper_prepare: this looper will accept calls to
	 * ALooper_addFd() that do not have a callback (that is provide NULL
	 * for the callback).  In this case the caller of ALooper_pollOnce()
	 * or ALooper_pollAll() MUST check the return from these functions to
	 * discover when data is available on such fds and process it.
	 */
	LOOPER_PREPARE_ALLOW_NON_CALLBACKS = C.ALOOPER_PREPARE_ALLOW_NON_CALLBACKS
)

/**
 * Prepares a looper associated with the calling thread, and returns it.
 * If the thread already has a looper, it is returned.  Otherwise, a new
 * one is created, associated with the thread, and returned.
 *
 * The opts may be ALOOPER_PREPARE_ALLOW_NON_CALLBACKS or 0.
 */
func looperPrepare(opts int) (looper *Looper) {
	return (*Looper)(C.ALooper_prepare(C.int(opts)))
}

const (
	/**
	 * Result from ALooper_pollOnce() and ALooper_pollAll():
	 * The poll was awoken using wake() before the timeout expired
	 * and no callbacks were executed and no other file descriptors were ready.
	 */
	LOOPER_POLL_WAKE = C.ALOOPER_POLL_WAKE

	/**
	 * Result from ALooper_pollOnce() and ALooper_pollAll():
	 * One or more callbacks were executed.
	 */
	LOOPER_POLL_CALLBACK = C.ALOOPER_POLL_CALLBACK

	/**
	 * Result from ALooper_pollOnce() and ALooper_pollAll():
	 * The timeout expired.
	 */
	LOOPER_POLL_TIMEOUT = C.ALOOPER_POLL_TIMEOUT

	/**
	 * Result from ALooper_pollOnce() and ALooper_pollAll():
	 * An error occurred.
	 */
	LOOPER_POLL_ERROR = C.ALOOPER_POLL_ERROR
)

type Looper C.ALooper

func (looper *Looper) cptr() *C.ALooper {
	return (*C.ALooper)(looper)
}

/**
 * Acquire a reference on the given ALooper object.  This prevents the object
 * from being deleted until the reference is removed.  This is only needed
 * to safely hand an ALooper from one thread to another.
 */
func (looper *Looper) Acquire() {
	C.ALooper_acquire(looper.cptr())
}

/**
 * Remove a reference that was previously acquired with ALooper_acquire().
 */
func (looper *Looper) Release() {
	C.ALooper_release(looper.cptr())
}

/**
 * Flags for file descriptor events that a looper can monitor.
 *
 * These flag bits can be combined to monitor multiple events at once.
 */

const (
	/**
	 * The file descriptor is available for read operations.
	 */
	LOOPER_EVENT_INPUT = C.ALOOPER_EVENT_INPUT

	/**
	 * The file descriptor is available for write operations.
	 */
	LOOPER_EVENT_OUTPUT = C.ALOOPER_EVENT_OUTPUT

	/**
	 * The file descriptor has encountered an error condition.
	 *
	 * The looper always sends notifications about errors; it is not necessary
	 * to specify this event flag in the requested event set.
	 */
	LOOPER_EVENT_ERROR = C.ALOOPER_EVENT_ERROR

	/**
	 * The file descriptor was hung up.
	 * For example, indicates that the remote end of a pipe or socket was closed.
	 *
	 * The looper always sends notifications about hangups; it is not necessary
	 * to specify this event flag in the requested event set.
	 */
	LOOPER_EVENT_HANGUP = C.ALOOPER_EVENT_HANGUP

	/**
	 * The file descriptor is invalid.
	 * For example, the file descriptor was closed prematurely.
	 *
	 * The looper always sends notifications about invalid file descriptors; it is not necessary
	 * to specify this event flag in the requested event set.
	 */
	LOOPER_EVENT_INVALID = C.ALOOPER_EVENT_INVALID
)

/**
 * For callback-based event loops, this is the prototype of the function
 * that is called when a file descriptor event occurs.
 * It is given the file descriptor it is associated with,
 * a bitmask of the poll events that were triggered (typically ALOOPER_EVENT_INPUT),
 * and the data pointer that was originally supplied.
 *
 * Implementations should return 1 to continue receiving callbacks, or 0
 * to have this file descriptor and callback unregistered from the looper.
 */
///typedef int (*ALooper_callbackFunc)(int fd, int events, void* data);
type LooperCallback func(fd, events int, data unsafe.Pointer) int

/**
 * Waits for events to be available, with optional timeout in milliseconds.
 * Invokes callbacks for all file descriptors on which an event occurred.
 *
 * If the timeout is zero, returns immediately without blocking.
 * If the timeout is negative, waits indefinitely until an event appears.
 *
 * Returns ALOOPER_POLL_WAKE if the poll was awoken using wake() before
 * the timeout expired and no callbacks were invoked and no other file
 * descriptors were ready.
 *
 * Returns ALOOPER_POLL_CALLBACK if one or more callbacks were invoked.
 *
 * Returns ALOOPER_POLL_TIMEOUT if there was no data before the given
 * timeout expired.
 *
 * Returns ALOOPER_POLL_ERROR if an error occurred.
 *
 * Returns a value >= 0 containing an identifier if its file descriptor has data
 * and it has no callback function (requiring the caller here to handle it).
 * In this (and only this) case outFd, outEvents and outData will contain the poll
 * events and data associated with the fd, otherwise they will be set to NULL.
 *
 * This method does not return until it has finished invoking the appropriate callbacks
 * for all file descriptors that were signalled.
 */
func looperPollOnce(timeoutMillis int) (ident, outFD, outEvents int, outData uintptr) {
	var coutFd, coutEvents C.int
	var coutData unsafe.Pointer
	cret := C.ALooper_pollOnce(C.int(timeoutMillis), &coutFd, &coutEvents, &coutData)
	return int(cret), int(coutFd), int(coutEvents), uintptr(coutData)
}

/**
 * Like ALooper_pollOnce(), but performs all pending callbacks until all
 * data has been consumed or a file descriptor is available with no callback.
 * This function will never return ALOOPER_POLL_CALLBACK.
 */
func looperPollAll(timeoutMillis int) (ident, outFD, outEvents int, outData uintptr) {
	var coutFd, coutEvents C.int
	var coutData unsafe.Pointer
	cident := C.ALooper_pollAll(C.int(timeoutMillis), &coutFd, &coutEvents, &coutData)
	return int(cident), int(coutFd), int(coutEvents), uintptr(coutData)
}

/**
 * Wakes the poll asynchronously.
 *
 * This method can be called on any thread.
 * This method returns immediately.
 */
func (looper *Looper) Wake() {
	C.ALooper_wake(looper.cptr())
}

/**
 * Adds a new file descriptor to be polled by the looper.
 * If the same file descriptor was previously added, it is replaced.
 *
 * "fd" is the file descriptor to be added.
 * "ident" is an identifier for this event, which is returned from ALooper_pollOnce().
 * The identifier must be >= 0, or ALOOPER_POLL_CALLBACK if providing a non-NULL callback.
 * "events" are the poll events to wake up on.  Typically this is ALOOPER_EVENT_INPUT.
 * "callback" is the function to call when there is an event on the file descriptor.
 * "data" is a private data pointer to supply to the callback.
 *
 * There are two main uses of this function:
 *
 * (1) If "callback" is non-NULL, then this function will be called when there is
 * data on the file descriptor.  It should execute any events it has pending,
 * appropriately reading from the file descriptor.  The 'ident' is ignored in this case.
 *
 * (2) If "callback" is NULL, the 'ident' will be returned by ALooper_pollOnce
 * when its file descriptor has data available, requiring the caller to take
 * care of processing it.
 *
 * Returns 1 if the file descriptor was added or -1 if an error occurred.
 *
 * This method can be called on any thread.
 * This method may block briefly if it needs to wake the poll.
 */
///int ALooper_addFd(ALooper* looper, int fd, int ident, int events,
///        ALooper_callbackFunc callback, void* data);
//export cgoLooperCallback
func cgoLooperCallback(fd, events C.int, data unsafe.Pointer) C.int {
	d := (*looperData)(data)
	return C.int(d.callback(int(fd), int(events), d.data))
}

type looperData struct {
	callback LooperCallback
	data     unsafe.Pointer
}

func (looper *Looper) AddFd(fd, ident, events int,
	callback LooperCallback, data unsafe.Pointer) int {
	if callback == nil {
		return int(C.ALooper_addFd(looper.cptr(), C.int(fd), C.int(ident), C.int(events),
			nil, data))
	} else {
		return int(C.ALooper_addFd(looper.cptr(), C.int(fd), C.int(ident), C.int(events),
			(*[0]byte)(C.cgoLooperCallback), unsafe.Pointer(&looperData{callback: callback, data: data})))
	}
}

/**
 * Removes a previously added file descriptor from the looper.
 *
 * When this method returns, it is safe to close the file descriptor since the looper
 * will no longer have a reference to it.  However, it is possible for the callback to
 * already be running or for it to run one last time if the file descriptor was already
 * signalled.  Calling code is responsible for ensuring that this case is safely handled.
 * For example, if the callback takes care of removing itself during its own execution either
 * by returning 0 or by calling this method, then it can be guaranteed to not be invoked
 * again at any later time unless registered anew.
 *
 * Returns 1 if the file descriptor was removed, 0 if none was previously registered
 * or -1 if an error occurred.
 *
 * This method can be called on any thread.
 * This method may block briefly if it needs to wake the poll.
 */
///int ALooper_removeFd(ALooper* looper, int fd);
func (looper *Looper) RemoveFd(fd int) int {
	return int(C.ALooper_removeFd(looper.cptr(), C.int(fd)))
}
