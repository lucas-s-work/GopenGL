package graphics

import (
	"gopengl/graphics/opengl"
	"unsafe"
)

/*
Go routine calls for graphics methods
*/

type RenderObjectJob struct {
	obj    *RenderObject
	name   byte
	params []interface{}
	retVal unsafe.Pointer
}

type VAOJob struct {
	obj    *opengl.VAO
	name   byte
	params []interface{}
	retVal unsafe.Pointer
}

//Job types
const (
	CREATE_RENDER_OBJECT byte = iota
	ADD_SQUARE           byte = iota
	MODIFY_VERT_SQUARE   byte = iota
	MODIFY_TEX_SQUARE    byte = iota
	TRANSLATE            byte = iota
	ROTATE               byte = iota
	MODIFY_VERT_RECT     byte = iota
	MODIFY_TEX_RECT      byte = iota
	ADD_RECT             byte = iota
)

//Job queues
var (
	RenderObjectQueue  = make(chan RenderObjectJob)
	VAOQueue           = make(chan VAOJob)
	RenderObjectJobMap = make(map[byte]func(RenderObjectJob))
	VAOJobMap          = make(map[byte]func(VAOJob))
)

var (
	alive = true
)

/*
Job handling
*/

func Listen() {
	// Setup job maps
	RenderObjectJobMap[CREATE_RENDER_OBJECT] = callCreateRenderObject
	RenderObjectJobMap[ADD_SQUARE] = callAddSquare
	RenderObjectJobMap[MODIFY_VERT_SQUARE] = callModifyVertSquare
	RenderObjectJobMap[MODIFY_TEX_SQUARE] = callModifyTexSquare
	RenderObjectJobMap[TRANSLATE] = callTranslate
	RenderObjectJobMap[ROTATE] = callRotate
	RenderObjectJobMap[MODIFY_VERT_RECT] = callModifyVertRect
	RenderObjectJobMap[MODIFY_TEX_RECT] = callModifyTexRect
	RenderObjectJobMap[ADD_RECT] = callAddRect

	defer cleanUp()

	for !ShouldClose() {
		select {
		case job := <-RenderObjectQueue:
			callRenderObjectJob(job)
		case job := <-VAOQueue:
			callVAOJob(job)
		default:
			Render()
		}
	}

	alive = false
}

var (
	maxCompletedJobs uint16 = 500
	completedJobs    uint16
)

func callRenderObjectJob(job RenderObjectJob) {
	RenderObjectJobMap[job.name](job)

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

func callTranslate(job RenderObjectJob) {
	params := job.params

	job.obj.Translate(
		params[0].(float32),
		params[1].(float32),
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

/*
Graphics job methods, these enqueue the job to be performed, graphics.go methods MUST NOT be used directly on RenderObjects generated here
These are all called *Outside* the main thread which the opengl context is running on.
*/

func CreateRenderObjectJob(ro *RenderObject, size int, texture string, defaultShader bool) {
	RenderObjectQueue <- RenderObjectJob{
		ro,
		CREATE_RENDER_OBJECT,
		[]interface{}{size, texture, defaultShader},
		nil,
	}
}

func (obj *RenderObject) AddSquareJob(x, y, xTex, yTex, width, widthTex float32) *int {
	freeVert := 0

	RenderObjectQueue <- RenderObjectJob{
		obj,
		ADD_SQUARE,
		[]interface{}{x, y, xTex, yTex, width, widthTex},
		unsafe.Pointer(&freeVert),
	}

	return &freeVert
}

func (obj *RenderObject) ModifyVertSquareJob(index *int, x, y, width float32) {
	RenderObjectQueue <- RenderObjectJob{
		obj,
		MODIFY_VERT_SQUARE,
		[]interface{}{*index, x, y, width},
		nil,
	}
}

func (obj *RenderObject) ModifyTexSquareJob(index *int, x, y, width float32) {
	RenderObjectQueue <- RenderObjectJob{
		obj,
		MODIFY_TEX_SQUARE,
		[]interface{}{*index, x, y, width},
		nil,
	}
}

func (obj *RenderObject) TranslateJob(x, y float32) {
	RenderObjectQueue <- RenderObjectJob{
		obj,
		TRANSLATE,
		[]interface{}{x, y},
		nil,
	}
}

func (obj *RenderObject) RotateJob(x, y, rot float32) {
	RenderObjectQueue <- RenderObjectJob{
		obj,
		ROTATE,
		[]interface{}{x, y, rot},
		nil,
	}
}

func (obj *RenderObject) AddRectJob(x, y, xTex, yTex, width, height, widthTex, heightTex float32) *int {
	freeVert := 0

	RenderObjectQueue <- RenderObjectJob{
		obj,
		ADD_RECT,
		[]interface{}{x, y, xTex, yTex, width, height, widthTex, heightTex},
		unsafe.Pointer(&freeVert),
	}

	return &freeVert
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
