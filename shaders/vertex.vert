#version 410
in vec2 vert;
in vec2 verttexcoord;

uniform vec2 trans;
uniform vec3 rot;

out vec2 fragtexcoord;
void main(){
    fragtexcoord=verttexcoord;
    
    //rotate before translate
    vec2 pos=vert;
    vec2 rotcenter=vec2(rot.x,rot.y);
    pos-=rotcenter;
    mat2 rotmat=mat2(cos(rot.z),-sin(rot.z),sin(rot.z),cos(rot.z));
    pos*=rotmat;
    pos+=rotcenter;
    pos+=trans;
    gl_Position=vec4(pos,0.,1.);
}