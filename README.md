# Gopengl
Gopengl's goal is to produce a high level, easily scalable opengl application in go. The core idea is to allow the use of go routines to generate opengl objects
with very few requirements to be cautious about race conditions.

Gopengl uses the go opengl bindings from: 

## Initialization
``` go 
// Thread dealing with the opengl context MUST be locked to the main thread
func init(){
    runtime.LockOSThread()
}

func main(){
    //Initialize gopengl
    graphics.Init()
    //create our window
    window := graphics.CreateWindow(800,600, "test application")
    //Assign the window to the opengl context
    graphics.SetWindow(window)
    grpahics.SetWindowSize(800,600)

    // ... do stuff ...
}
```
Window hints are currently unsupported.

## Single threaded functions
When running only in the main thread the use of functions in the `graphics.go` file can be used.
Due to Gopengl being used in Battleships functions will only be created as they are needed, all current ones rely on rectangles
however extra support will be added in the future. If required, lower level calls can be made directly to the VAO's but these are not
multithreadable.

After initialization these methods can be called.

### Render objects
At the high level render objects are what everything is called through. Render objects are an abstraction on top of a VAO, they are designed to be used for
objects that move and rotate together, for instance a character with clothes, the player and clothes should move and rotate together but we may have multuple players.

All render objects are handled by a central controller that handles cleanup, rendering & filtering. In future setting camera position will be supported 
globally for all render objects.

#### Creating a render object
``` go
func main(){
    // .. initialize ...

    var ro graphics.RenderObject
    graphics.CreateRenderObject(&ro, vertNum int, texturePath string, defaultShader bool)
}
```
#### Adding a square and rectangle
``` go
// Create the square
square := ro.AddSquare(x, y int, xTex, yTex int, width int, texWidth int)

// Modify the square
ro.ModifySquare(square, x, y, xTex, yTex, width, texWidth)
```

If modifying textures or vertices only then there exists `ro.ModifySquareVert` and `ro.ModifySquareTex`.

#### Transforming render objects
You can peform transformations and rotations on every square/object contained within a render object

```go
// Rotate the renderobject about centre x, y by rotation radians
ro.Rotate(x,y, rotation float32)
// Translate the render object by x, y
ro.Translate(x,y float32)
```

## Multi theaded functions
Multithreading graphics calls is performed by enqueuing jobs instead of performing them immediately. The main go routine of your application becomes 
solely dedicated to processing these graphics calls and all other go routines are performed elsewhere. Note that it is currently not possible to have
multiple opengl contexts run like this, support is unlikely to ever be added.

### Listening
Jobs are listened for and executed just by calling `graphics.Listen()`
```go
func main(){
    // initialize as normal

    // startup the rest of your program
    go restofapplication()

    // Start listening for jobs
    graphics.Listen()
}
```

### Job execution
Jobs are named analagously to the original function by adding a `Job` suffix, for instance creating a render object.
```go
var ro graphics.RenderObject
graphics.CreateRenderObjectJob(&ro, vertNum int, texturePath string, defaultShader bool)
```

While there is no guarantee when jobs will be performed the order of the jobs can be guaranteed so multithreaded code can be written synchronously. 

In order to avoid any possible issues with this avoid using a single render object across multiple go routines.

## VAO's
There is currently support for direct VAO interaction for single threaded uses, once required in Battleships multithreaded support will be added.

VAO's have similar grouped rotation and translation functions but are not automatically rendered or cleaned up.

### VAO Creation
``` go
vao := opengl.CreateVAO(vertNum int, textureSource string, defaultShader bool)
```
### Buffer handling
Multiple methods are used to create and handle buffers, whilst the VAO object does check if it has initialized it's buffers on calling `UpdateBuffers`, it is still good for clarity to directly call `CreateBuffers`.

### Setting buffer data
VAO's handle two copies of the buffer data, a buffer stored in the VAO struct and the data stored directly on the gpu. We can set both at once or modify the struct data and later modify the gpu data.

#### Setting vao struct data
`SetData` Updates the VAO struct buffers but does not modify the GPU buffers.
```go
vao.SetData(vertData, texData float32[])
```
#### Creating and Updating Buffers
`CreateBuffers` and `UpdateBuffers` respectively initialize or update the GPU's buffer data with the current VAO buffer.
``` go
// Create buffers with existing buffer data
vao.CreateBuffers()

// Set new buffer data and update GPU buffers
vao.SetData(newVertData, newTexData)
vao.UpdateBuffers()

// For shorthand update both at the same time
vao.UpdateBufferData(newVertData, newTexData)
```

Analagous functions `UpdateVertBuffer` & `UpdateTexBuffer` for updating specific vertData and texData

Often you don't need to update entire buffers at once, for this use `UpdateBufferIndex`.
```go
vao.UpdateBufferIndex(floatIndex, vertData, texData)
```

Analgous functions `UpdateVertBufferIndex` & `UpdateTexBufferIndex` for updating specifically vertData or texData.

### VAO Deletion
```go
vao.Delete()
```
### VAO Rendering
Rendering VAO's is abstracted using `PrepRender` & `FinishRender` which occur before and after rendering, these functions Bind the vao and buffers and setup the textures and shaders.

```go
vao.PrepRender()
gl.DrawArrays(gl.TRIANGLES, 0, vertNum)
vao.FinishRender()
```

A shorthand exists for this which allows multiple VAO's to render at once.

```go
vaos := make([]*opengl.VAO, 1)
vaos[0] = vao

opengl.Render(vaos)
```

## Shaders
The `defaultShader` option used when creating VAO's and RenderObjects determines if on creation the basic shaders should be supported. If using custom shaders `defaultShader` should be false, also note that the `Translate` & `Rotate` methods will not work.

### Loading custom shaders
If you already have an created ShaderProgram you can pass it's id as an arguement to `CreateProgram`.
``` go 
// Create a new shader program
program := opengl.CreateProgram(0)

// Wrap an existing shader program
program := opengl.CreateProgram(programId)
```

Shaders are loaded relative to the current working directory, currently only loading and variable management exist for vertex and fragment shaders.
```go
program.LoadVertShader(relPath)
program.LoadFragShader(relPath)
```

Parsing and compilation of all shader types is supported however.
```go
program := opengl.CreateProgram(0)
shaderData, _ := ioutil.ReadFile(shaderPath)

// Load shader, SHADERTYPE is a Go-opengl uint32 value
// eg gl.VERTEX_SHADER, gl.GEOMETRY_SHADER
program.LoadShader(shaderData, gl.SHADERTYPE)
```

### Linking and Using shaders
Once all required shaders are loaded call `program.Link()` to link the shaders together. This does not need to be called each render only when changing a programs current shaders.

To use a shader call `program.Use()`
### Attributes and Uniforms
**This section is likely to change**

#### Attributes
Due to how attributes are set in opengl the current implementation requires these to be enabled on VAO binding.

Also note due to this if using a custom shader `vao.CreateBuffers()` cannot be used, manual vao buffer creation will need to be perfomed.

```go
func SetupVAO(){
    // ... do stuff ...
    // Attach and attribute to this program
    program.AddAttribute("attributeName")
    // ... do stuff ...

    BindVAOFunction()
    program.EnableAttribute("attributeName")
}
```

#### Uniforms
Uniform values can be changed anywhere in code and are simpler to use because of this.

```go
func SetupVAO(){
    // ... do stuff ...
    // Add a uniform value
    program.AddUniform("uniformName", uniformValue interface{})
    // ... do stuff ...
    // Modify an existing uniform value
    program.SetUniform("uniformName", uniformValue interface{})
}
```

Currently this uniform interface is just a wrapper for opengl methods, in future modifying uniforms and attributes will be performed like so.

``` go
// Create variables and attach them to the shader
uni := program.CreateUniform("uniformName")
attribute := program.CreateAttribute("attributeName")

// Set variable values without needing the program
attribute.SetData(attributeData)
uni.SetValue(value interface{})
```

# Planned
Features will be added as required by the Battleships project, some currently planned features are:
 - [] Triangle support in render objects
 - [] Simplified attribute and unnfirom usage
 - [] geometry shader support
 - [] more vao flexibility
 - [] window hints
 - [] better error reporting
 - [] full input handling
 - [] window poll jobs
 - [] prioritised jobs
 - [] vao multithreading support
 - [] logging

Currently all error handling is to outright panic, this is because most opengl errors will simply cause a sigterm so error reporting can be difficult, in future catching of errors will be reported especially around jobs with the option to not fail on certain jobs.
