type worker struct
{
	ip string
	port string
}

workerstatus = map[worker]status

func main() {
	manageWorkers()
	listenJob()
}

func manageWorkers() {
	listen from workers : 8080
	for accept {
		if unknown worker {
			new worker
		}
		if add {
			append(workers, worker)
		}
		else if remove {
			remove(workers, worker)
		}
	}
}

func listenJob() {
	jobs = []
	listen from client : 8081
	for accept {
		if add {
			jobs.add(job)
			go manageJob(job)
		}
		else if stop {
			jobs.stop(job)
		}
		elif remove {
			jobs.remove(job)
		}
	}
}


func manageJob() {
	divide job to tasks
	for {
		if not tasks.remain() {
			return
		}
		findWorker w
		taskInfo info
		postTask(w, info)
	}
}

func postTask(w, info) {
	connect to worker
	send task info
	change task status to 'processing'
	waiting job done
	if err {
		add task error
		change task status to 'waiting'
	}
	else {
		change task status to 'done'
		return
	}
}
