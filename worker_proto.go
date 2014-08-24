func register() {
	wakeup
	register to lazyfarm : 8080
}

func waiting() {
	listen from lazyfarm
	for accept {
		do work
	}
}

