#! /bin/bash

classname="$(basename -s .go $1)"
classinit=${classname:0:1}
parent_dir_name=$(basename $(dirname $1))

cat >> $1 <<EOF
package $parent_dir_name

import "github.com/prizelobby/ebitengine-template/ui"

type $2 struct {

}

func New$2() *$2 {
    return &$2{}
}

func ($classinit *$2) Update() {

}

func ($classinit *$2) Draw(screen *ui.ScaledScreen) {
    
}
EOF