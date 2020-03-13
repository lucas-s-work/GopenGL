#version 410
in vec2 vert;
in vec4 rotgroup;
in vec2 verttexcoord;

//Translation, window dimension scaling, rotation
uniform vec2 trans;
uniform vec2 dim;
uniform vec4 rot;

out vec2 fragtexcoord;
void main(){
    // Set tex coords for frag shader
    fragtexcoord=verttexcoord;
    vec2 pos=vert;
    
    //Apply rotgroup rotation first, we want local changes then global changes to each vertex
    vec2 rotcenter=vec2(rotgroup.x,rotgroup.y);
    pos-=rotcenter;
    
    mat2 rotmat=mat2(
        rotgroup.z,rotgroup.w,
        -rotgroup.w,rotgroup.z
    );
    pos=rotmat*pos;
    
    pos+=rotcenter;
    
    // Apply uniform rotation
    rotcenter=vec2(rot.x,rot.y);
    pos=pos-rotcenter;
    
    rotmat=mat2(
        rot.z,rot.w,
        -rot.w,rot.z
    );
    
    pos=rotmat*pos;
    
    pos=pos+rotcenter;
    
    // Apply screen scaling from pixel coordinates
    pos.x=(pos.x/(.5*dim.x))-1;
    pos.y=1-(pos.y/(.5*dim.y));
    
    gl_Position=vec4(pos+trans,0.,1.);
}