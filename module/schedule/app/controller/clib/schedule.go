package clib

import (
	"fmt"
	"net/http"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/views/vadmin"
)

const scheduleBC = "Schedule||/admin/schedule**stopwatch"

func ScheduleList(w http.ResponseWriter, r *http.Request) {
	controller.Act("schedule.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		jobs := as.Services.Schedule.ListJobs()
		ps.SetTitleAndData("Schedules", jobs)
		ecs := as.Services.Schedule.ExecCounts
		return controller.Render(r, as, &vadmin.Schedule{Jobs: jobs, ExecCounts: ecs}, ps, "admin", scheduleBC)
	})
}

func ScheduleDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("schedule.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		id, err := cutil.PathUUID(r, "id")
		if err != nil {
			return "", err
		}

		job := as.Services.Schedule.GetJob(*id)
		if job == nil {
			msg := fmt.Sprintf("no scheduled job with id [%s]", id)
			return controller.FlashAndRedir(false, msg, "/admin/schedule", ps)
		}
		res := as.Services.Schedule.Results[*id]
		ec := as.Services.Schedule.ExecCounts[*id]

		ps.SetTitleAndData(job.ID.String(), job)
		page := &vadmin.ScheduleDetail{Job: job, Result: res, ExecCount: ec}
		return controller.Render(r, as, page, ps, "admin", scheduleBC, job.ID.String())
	})
}
