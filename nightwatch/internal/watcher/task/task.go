package task

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/jianghushinian/blog-go-example/nightwatch/pkg/model"
)

func toJob(task *model.Task) *batchv1.Job {
	backoffLimit := int32(1)
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: task.Name,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "task-job",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "task",
							Image:   task.Info.Image,
							Command: task.Info.Command,
							Args:    task.Info.Args,
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backoffLimit,
		},
	}
	return jobSpec
}

func toTaskStatus(job *batchv1.Job) model.TaskStatus {
	// 检查 Job 状态
	switch {
	case isJobCompleted(job):
		return model.TaskStatusSucceeded
	case isJobFailed(job):
		return model.TaskStatusFailed
	case isJobSuspended(job):
		return model.TaskStatusPending
	case isJobFailureTarget(job):
		return model.TaskStatusFailed
	case isJobSuccessCriteriaMet(job):
		return model.TaskStatusSucceeded
	case isJobActive(job):
		return model.TaskStatusRunning
	default:
		return model.TaskStatusUnknown
	}
}

// 判断 Job 是否完成
func isJobCompleted(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// 判断 Job 是否失败
func isJobFailed(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobFailed && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// 判断 Job 是否挂起
func isJobSuspended(job *batchv1.Job) bool {
	return job.Spec.Suspend != nil && *job.Spec.Suspend
}

// 判断 Job 是否即将失败
func isJobFailureTarget(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobFailureTarget && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// 判断 Job 成功条件是否已满足
func isJobSuccessCriteriaMet(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobSuccessCriteriaMet && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// 判断 Job 是否正在运行
func isJobActive(job *batchv1.Job) bool {
	return job.Status.Active > 0
}
