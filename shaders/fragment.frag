#version 410
uniform sampler2D tex;

out vec4 frag_colour;
in vec2 fragtexcoord;
void main(){
    frag_colour=texture(tex, fragtexcoord);
}