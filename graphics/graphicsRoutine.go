package graphics

import (
	"gopengl/graphics/opengl"
	"time"
	"unsafe"
)

/*
Go routine calls for graphics methods
*/

type RenderObjectJob struct {
	obj      *RenderObject
	params   []interface{}
	retVal   unsafe.Pointer
	callable func(RenderObjectJob)
}

type VAOJob struct {
	obj    *opengl.VAO
	name   byte
	params []interface{}
	retVal unsafe.Pointer
}

/*
TODO
Remove job mappings, pass function as pointer into job instead.
*/

//Job queues
var (
	RenderObjectQueue = make(chan RenderObjectJob)
	VAOQueue          = make(chan VAOJob)
	VAOJobMap         = make(map[byte]func(VAOJob))
)

var (
	alive = true
)

/*
Job handling
*/

func Listen() {
	if &renderDelta == nil {
		SetFrameRate(60, 5)
	}

	defer cleanUp()

	for !ShouldClose() {
		select {
		case job := <-RenderObjectQueue:
			callRenderObjectJob(job)
		case job := <-VAOQueue:
			callVAOJob(job)
		default:
			t := time.Now()

			if t.Sub(lastRender).Nanoseconds() >= renderDelta {
				lastRender = t
				Render()
			} else {
				// time.Sleep(renderSleep)
			}

		}
	}

	alive = false
}

// SetFrameRate ... Rate: Frame Rate in fps, Sampling: The maximum number of times to check the render queue inbetween frames.
func SetFrameRate(rate int, sampling int) {
	renderDelta = int64(1000000000 / rate)
	renderSleep = time.Duration(1000000000 / (sampling * rate))
}

var (
	maxCompletedJobs uint16 = 500
	completedJobs    uint16
	lastRender       time.Time
	renderDelta      int64
	renderSleep      time.Duration
)

func callRenderObjectJob(job RenderObjectJob) {
	job.callable(job)

	checkRender()
}

func callVAOJob(job VAOJob) {
	VAOJobMap[job.name](job)

	checkRender()
}

func checkRender() {
	completedJobs++

	if completedJobs >= maxCompletedJobs {
		Render()
		completedJobs = 0
	}
}

/*
RenderObjectJob name mappings, these must not return any values instead they must be passed pointers to be modified.
These are all called *Inside* the main thread which the opengl context is running on.
*/

func callCreateRenderObject(job RenderObjectJob) {
	CreateRenderObject(job.obj, job.params[0].(int), job.params[1].(string), job.params[2].(bool))
}

func callAddSquare(job RenderObjectJob) {
	params := job.params

	freeVert := job.obj.AddSquare(
		params[0].(float32),
		params[1].(float32),
		params[2].(float32),
		params[3].(float32),
		params[4].(float32),
		params[5].(float32),
	)

	*(*int)(job.retVal) = freeVert
}

func callModifyVertSquare(job RenderObjectJob) {
	params := job.params

	job.obj.ModifyVertSquare(
		params[0].(int),
		params[1].(float32),
		params[2].(float32),
		params[3].(float32),
	)
}

func callModifyTexSquare(job RenderObjectJob) {
	params := job.params

	job.obj.ModifyTexSquare(
		params[0].(int),
		params[1].(float32),
		params[2].(float32),
		params[3].(float32),
	)
}

func callRotate(job RenderObjectJob) {
	params := job.params

	job.obj.Rotate(
		params[0].(float32),
		params[1].(float32),
		params[2].(float32),
	)
}

func callAddRect(job RenderObjectJob) {
	params := job.params

	freeVert := job.obj.AddRect(
		params[0].(float32),
		params[1].(float32),
		params[2].(float32),
		params[3].(float32),
		params[4].(float32),
		params[5].(float32),
		params[6].(float32),
		params[7].(float32),
	)

	*(*int)(job.retVal) = freeVert
}

func callModifyTexRect(job RenderObjectJob) {
	params := job.params

	job.obj.ModifyTexRect(
		params[0].(int),
		params[1].(float32),
		params[2].(float32),
		params[3].(float32),
		params[4].(float32),
	)
}

func callModifyVertRect(job RenderObjectJob) {
	params := job.params

	job.obj.ModifyVertRect(
		params[0].(int),
		params[1].(float32),
		params[2].(float32),
		params[3].(float32),
		params[4].(float32),
	)
}

func callResetGroupedRotation(job RenderObjectJob) {
	job.obj.ResetGroupedRotation()
}

func callSetAllGroupedRotation(job RenderObjectJob) {
	params := job.params

	job.obj.SetAllGroupedRotation(
		params[0].(float32),
		params[1].(float32),
		params[2].(float32),
	)
}

func callSetGroupedRotation(job RenderObjectJob) {
	params := job.params

	job.obj.SetGroupedRotation(
		params[0].(float32),
		params[1].(float32),
		params[2].(float32),
		params[3].(int),
		params[4].(int),
	)
}

func callUpdateBuffers(job RenderObjectJob) {
	job.obj.vao.UpdateBuffers()
}

/*
Graphics job methods, these enqueue the job to be performed, graphics.go methods MUST NOT be used directly on RenderObjects generated here
These are all called *Outside* the main thread which the opengl context is running on.
*/

func CreateRenderObjectJob(ro *RenderObject, size int, texture string, defaultShader bool) {
	RenderObjectQueue <- RenderObjectJob{
		ro,
		[]interface{}{size, texture, defaultShader},
		nil,
		callCreateRenderObject,
	}
}

func (obj *RenderObject) AddSquareJob(x, y, xTex, yTex, width, widthTex float32) *int {
	freeVert := 0

	RenderObjectQueue <- RenderObjectJob{
		obj,
		[]interface{}{x, y, xTex, yTex, width, widthTex},
		unsafe.Pointer(&freeVert),
		callAddSquare,
	}

	return &freeVert
}

func (obj *RenderObject) ModifyVertSquareJob(index *int, x, y, width float32) {
	RenderObjectQueue <- RenderObjectJob{
		obj,
		[]interface{}{*index, x, y, width},
		nil,
		callModifyVertSquare,
	}
}

func (obj *RenderObject) ModifyTexSquareJob(index *int, x, y, width float32) {
	RenderObjectQueue <- RenderObjectJob{
		obj,
		[]interface{}{*index, x, y, width},
		nil,
		callModifyTexSquare,
	}
}

func (obj *RenderObject) RotateJob(x, y, rot float32) {
	RenderObjectQueue <- RenderObjectJob{
		obj,
		[]interface{}{x, y, rot},
		nil,
		callRotate,
	}
}

func (obj *RenderObject) ResetGroupedRotationJob() {
	RenderObjectQueue <- RenderObjectJob{
		obj,
		nil,
		nil,
		callResetGroupedRotation,
	}
}

func (obj *RenderObject) SetAllGroupedRotationJob(x, y, rad float32) {
	RenderObjectQueue <- RenderObjectJob{
		obj,
		[]interface{}{x, y, rad},
		nil,
		callSetAllGroupedRotation,
	}
}

func (obj *RenderObject) SetGroupedRotationJob(x, y, rad float32, start, end int) {
	RenderObjectQueue <- RenderObjectJob{
		obj,
		[]interface{}{x, y, rad, start, end},
		nil,
		callSetGroupedRotation,
	}
}

func (obj *RenderObject) AddRectJob(x, y, xTex, yTex, width, height, widthTex, heightTex float32) *int {
	freeVert := 0

	RenderObjectQueue <- RenderObjectJob{
		obj,
		[]interface{}{x, y, xTex, yTex, width, height, widthTex, heightTex},
		unsafe.Pointer(&freeVert),
		callAddRect,
	}

	return &freeVert
}

func (obj *RenderObject) UpdateBuffersJob() {
	RenderObjectQueue <- RenderObjectJob{
		obj,
		nil,
		nil,
		callUpdateBuffers,
	}
}

/*
Cleanup
*/

func cleanUp() {
	DeleteRenderObjects()
	window.Destroy()
}

func ShouldClose() bool {
	return window.ShouldClose()
}

func Alive() *bool {
	return &alive
}

/*
Utility methods
*/

// CreateEmptyRenderObject ... create an empty render object for use with graphicsRoutine
func CreateEmptyRenderObject() *RenderObject {
	return &RenderObject{}
}

// CreateEmptyVAO ... create an empty vao object for use with graphicsRoutine
func CreateEmptyVAO() *opengl.VAO {
	return &opengl.VAO{}
}
