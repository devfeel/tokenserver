package task

import (
	"github.com/devfeel/tokenserver/framework/log"
	"sync"
	"time"
)

const (
	taskLogTarget  = "task"
	taskState_Init = "0"
	taskState_Run  = "1"
	taskState_Stop = "2"
)

//task 容器
type TaskService struct {
	taskMap         map[string]*TaskInfo
	taskMutex       *sync.RWMutex
	logPath         string
	taskHandleMap   map[string]TaskHandle
	taskHandleMutex *sync.RWMutex
}

//Task定义
type TaskInfo struct {
	TimeTicker *time.Ticker
	TaskID     string
	handler    TaskHandle
	Context    *TaskContext
	state      string //匹配 taskState_Init、taskState_Run、taskState_Stop
}

//Task上下文信息
type TaskContext struct {
	TaskID     string
	Interval   int64 //运行间隔时间，单位毫秒
	HandleName string
	TaskData   interface{}
}

func StartNewService() *TaskService {
	service := new(TaskService)
	service.taskMutex = new(sync.RWMutex)
	service.taskHandleMutex = new(sync.RWMutex)
	service.taskHandleMap = make(map[string]TaskHandle)
	service.taskMap = make(map[string]*TaskInfo)

	return service
}

type TaskHandle func(*TaskContext)

//get handler by name
func (service *TaskService) GetHandler(name string) (TaskHandle, bool) {
	service.taskHandleMutex.RLock()
	defer service.taskHandleMutex.RUnlock()
	handler, exists := service.taskHandleMap[name]
	return handler, exists
}

//set task handler with handlername
func (service *TaskService) SetHandler(handlerName string, handler TaskHandle) {
	service.taskHandleMap[handlerName] = handler
}

//创建Task对象
func (service *TaskService) CreateTask(taskID string, interval int64, handlerName string, taskData interface{}) *TaskInfo {
	context := new(TaskContext)
	context.HandleName = handlerName
	context.TaskID = taskID
	context.Interval = interval
	context.TaskData = taskData

	handler, exists := service.GetHandler(handlerName)
	if !exists {
		return nil
	}

	task := new(TaskInfo)
	task.TaskID = context.TaskID
	task.handler = handler
	task.state = taskState_Init
	task.Context = context

	service.taskMutex.Lock()
	service.taskMap[context.TaskID] = task
	service.taskMutex.Unlock()
	return task
}

//start timeticker
func (task *TaskInfo) Start() {
	if task.state == taskState_Init || task.state == taskState_Stop {
		task.state = taskState_Run
		task.TimeTicker = time.NewTicker(time.Duration(task.Context.Interval) * time.Millisecond)
		go func() {
			for {
				select {
				case <-task.TimeTicker.C:
					//TODO:do log
					task.handler(task.Context)
				}
			}
		}()
	}
}

//stop timeticker
func (task *TaskInfo) Stop() {
	if task.state == taskState_Stop {
		task.TimeTicker.Stop()
		task.state = taskState_Stop
	}
}

//remove all task
func (service *TaskService) RemoveAllTask() {
	logger.Info("Task::resetAllTask begin...", taskLogTarget)
	service.StopAllTask()
	service.taskMap = make(map[string]*TaskInfo)
}

//结束所有Task
func (service *TaskService) StopAllTask() {
	logger.Info("Task::StopAllTask begin...", taskLogTarget)
	for k, v := range service.taskMap {
		logger.Info("Task::StopAllTask => "+k, taskLogTarget)
		v.Stop()
	}
	logger.Info("Task::StopAllTask end["+string(len(service.taskMap))+"]", taskLogTarget)
}

//启动所有Task
func (service *TaskService) StartAllTask() {
	logger.Info("Task::StartAllTask begin...", taskLogTarget)
	for _, v := range service.taskMap {
		logger.Info("Task::StartAllTask::StartTask => "+v.TaskID, taskLogTarget)
		v.Start()
	}
	logger.Info("Task::StartAllTask end", taskLogTarget)
}
