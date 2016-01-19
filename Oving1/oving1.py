from threading import Thread
from threading import Lock


i=0
lock = Lock()

def thread1():
	global i
	for j in range(0,99999):
		lock.acquire()
		i += 1
		lock.release()

def thread2():
	global i
	for j in range(0,99998):
		lock.acquire()
		i -= 1
		lock.release()

def main():
	mainThread_1 = Thread(target = thread1, args = (),)
	mainThread_1.start()

	mainThread_2 = Thread(target = thread2,args=(),)
	mainThread_2.start()

	mainThread_1.join()
	mainThread_2.join()

	print(i)


main()
