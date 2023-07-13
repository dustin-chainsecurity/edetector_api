package channel

import (

)

var SignalChannel = make(chan string)
var TaskChangeChannel = make(chan [][]string, 2)