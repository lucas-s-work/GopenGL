#version 410
in vec2 vert;
in vec3 rotgroup;
in vec2 verttexcoord;

//Translation, window dimension scaling, rotation
uniform vec2 trans;
uniform vec2 dim;
uniform vec4 rot;

out vec2 fragtexcoord;
void main(){
    // Set tex coords for frag shader
    fragtexcoord=verttexcoord;
    
    // Translate
    vec2 pos=vert+trans;
    
    // Apply uniform rotation
    vec2 rotcenter=vec2(rot.x,rot.y);
    pos=pos-rotcenter;
    
    mat2 rotmat=mat2(
        rot.z,rot.w,
        -rot.w,rot.z
    );
    
    pos=rotmat*pos;
    
    pos=pos+rotcenter;
    
    // Apply rotgroup rotation
    // Not yet implemented
    
    // Apply screen scaling from pixel coordinates
    pos.x=(pos.x/(.5*dim.x))-1;
    pos.y=1-(pos.y/(.5*dim.y));
    
    gl_Position=vec4(pos,0.,1.);
}