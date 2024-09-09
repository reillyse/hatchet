// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: workflows.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addStepParents = `-- name: AddStepParents :exec
INSERT INTO "_StepOrder" ("A", "B")
SELECT
    step."id",
    $1::uuid
FROM
    unnest($2::text[]) AS parent_readable_id
JOIN
    "Step" AS step ON step."readableId" = parent_readable_id AND step."jobId" = $3::uuid
`

type AddStepParentsParams struct {
	ID      pgtype.UUID `json:"id"`
	Parents []string    `json:"parents"`
	Jobid   pgtype.UUID `json:"jobid"`
}

func (q *Queries) AddStepParents(ctx context.Context, db DBTX, arg AddStepParentsParams) error {
	_, err := db.Exec(ctx, addStepParents, arg.ID, arg.Parents, arg.Jobid)
	return err
}

const addWorkflowTag = `-- name: AddWorkflowTag :exec
INSERT INTO "_WorkflowToWorkflowTag" ("A", "B")
SELECT $1::uuid, $2::uuid
ON CONFLICT DO NOTHING
`

type AddWorkflowTagParams struct {
	ID   pgtype.UUID `json:"id"`
	Tags pgtype.UUID `json:"tags"`
}

func (q *Queries) AddWorkflowTag(ctx context.Context, db DBTX, arg AddWorkflowTagParams) error {
	_, err := db.Exec(ctx, addWorkflowTag, arg.ID, arg.Tags)
	return err
}

const countRoundRobinGroupKeys = `-- name: CountRoundRobinGroupKeys :one
SELECT
    COUNT(DISTINCT "concurrencyGroupId") AS total
FROM
    "WorkflowRun" r1
JOIN
    "WorkflowVersion" workflowVersion ON r1."workflowVersionId" = workflowVersion."id"
WHERE
    r1."tenantId" = $1::uuid AND
    workflowVersion."deletedAt" IS NULL AND
    r1."deletedAt" IS NULL AND
    (
        $2::"WorkflowRunStatus" IS NULL OR
        r1."status" = $2::"WorkflowRunStatus"
    ) AND
    workflowVersion."workflowId" = $3::uuid
`

type CountRoundRobinGroupKeysParams struct {
	Tenantid   pgtype.UUID           `json:"tenantid"`
	Status     NullWorkflowRunStatus `json:"status"`
	Workflowid pgtype.UUID           `json:"workflowid"`
}

func (q *Queries) CountRoundRobinGroupKeys(ctx context.Context, db DBTX, arg CountRoundRobinGroupKeysParams) (int64, error) {
	row := db.QueryRow(ctx, countRoundRobinGroupKeys, arg.Tenantid, arg.Status, arg.Workflowid)
	var total int64
	err := row.Scan(&total)
	return total, err
}

const countWorkflowRunsRoundRobin = `-- name: CountWorkflowRunsRoundRobin :one
SELECT COUNT(*) AS total
FROM
    "WorkflowRun" r1
JOIN
    "WorkflowVersion" workflowVersion ON r1."workflowVersionId" = workflowVersion."id"
WHERE
    r1."tenantId" = $1::uuid AND
    workflowVersion."deletedAt" IS NULL AND
    r1."deletedAt" IS NULL AND
    (
        $2::"WorkflowRunStatus" IS NULL OR
        r1."status" = $2::"WorkflowRunStatus"
    ) AND
    workflowVersion."workflowId" = $3::uuid AND
    r1."concurrencyGroupId" IS NOT NULL AND
    (
        $4::text IS NULL OR
        r1."concurrencyGroupId" = $4::text
    )
`

type CountWorkflowRunsRoundRobinParams struct {
	Tenantid   pgtype.UUID           `json:"tenantid"`
	Status     NullWorkflowRunStatus `json:"status"`
	Workflowid pgtype.UUID           `json:"workflowid"`
	GroupKey   pgtype.Text           `json:"groupKey"`
}

func (q *Queries) CountWorkflowRunsRoundRobin(ctx context.Context, db DBTX, arg CountWorkflowRunsRoundRobinParams) (int64, error) {
	row := db.QueryRow(ctx, countWorkflowRunsRoundRobin,
		arg.Tenantid,
		arg.Status,
		arg.Workflowid,
		arg.GroupKey,
	)
	var total int64
	err := row.Scan(&total)
	return total, err
}

const countWorkflows = `-- name: CountWorkflows :one
SELECT
    count(workflows) OVER() AS total
FROM
    "Workflow" as workflows
WHERE
    workflows."tenantId" = $1 AND
    workflows."deletedAt" IS NULL AND
    (
        $2::text IS NULL OR
        workflows."id" IN (
            SELECT
                DISTINCT ON(t1."workflowId") t1."workflowId"
            FROM
                "WorkflowVersion" AS t1
                LEFT JOIN "WorkflowTriggers" AS j2 ON j2."workflowVersionId" = t1."id"
            WHERE
                (
                    j2."id" IN (
                        SELECT
                            t3."parentId"
                        FROM
                            "public"."WorkflowTriggerEventRef" AS t3
                        WHERE
                            t3."eventKey" = $2::text
                            AND t3."parentId" IS NOT NULL
                    )
                    AND j2."id" IS NOT NULL
                    AND t1."workflowId" IS NOT NULL
                )
            ORDER BY
                t1."workflowId" DESC, t1."order" DESC
        )
    )
`

type CountWorkflowsParams struct {
	TenantId pgtype.UUID `json:"tenantId"`
	EventKey pgtype.Text `json:"eventKey"`
}

func (q *Queries) CountWorkflows(ctx context.Context, db DBTX, arg CountWorkflowsParams) (int64, error) {
	row := db.QueryRow(ctx, countWorkflows, arg.TenantId, arg.EventKey)
	var total int64
	err := row.Scan(&total)
	return total, err
}

const createJob = `-- name: CreateJob :one
INSERT INTO "Job" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "tenantId",
    "workflowVersionId",
    "name",
    "description",
    "timeout",
    "kind"
) VALUES (
    $1::uuid,
    coalesce($2::timestamp, CURRENT_TIMESTAMP),
    coalesce($3::timestamp, CURRENT_TIMESTAMP),
    $4::timestamp,
    $5::uuid,
    $6::uuid,
    $7::text,
    $8::text,
    $9::text,
    coalesce($10::"JobKind", 'DEFAULT')
) RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "workflowVersionId", name, description, timeout, kind
`

type CreateJobParams struct {
	ID                pgtype.UUID      `json:"id"`
	CreatedAt         pgtype.Timestamp `json:"createdAt"`
	UpdatedAt         pgtype.Timestamp `json:"updatedAt"`
	Deletedat         pgtype.Timestamp `json:"deletedat"`
	Tenantid          pgtype.UUID      `json:"tenantid"`
	Workflowversionid pgtype.UUID      `json:"workflowversionid"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	Timeout           string           `json:"timeout"`
	Kind              NullJobKind      `json:"kind"`
}

func (q *Queries) CreateJob(ctx context.Context, db DBTX, arg CreateJobParams) (*Job, error) {
	row := db.QueryRow(ctx, createJob,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Deletedat,
		arg.Tenantid,
		arg.Workflowversionid,
		arg.Name,
		arg.Description,
		arg.Timeout,
		arg.Kind,
	)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.WorkflowVersionId,
		&i.Name,
		&i.Description,
		&i.Timeout,
		&i.Kind,
	)
	return &i, err
}

const createSchedules = `-- name: CreateSchedules :many
INSERT INTO "WorkflowTriggerScheduledRef" (
    "id",
    "parentId",
    "triggerAt",
    "input"
) VALUES (
    gen_random_uuid(),
    $1::uuid,
    unnest($2::timestamp[]),
    $3::jsonb
) RETURNING id, "parentId", "triggerAt", "tickerId", input, "childIndex", "childKey", "parentStepRunId", "parentWorkflowRunId"
`

type CreateSchedulesParams struct {
	Workflowrunid pgtype.UUID        `json:"workflowrunid"`
	Triggertimes  []pgtype.Timestamp `json:"triggertimes"`
	Input         []byte             `json:"input"`
}

func (q *Queries) CreateSchedules(ctx context.Context, db DBTX, arg CreateSchedulesParams) ([]*WorkflowTriggerScheduledRef, error) {
	rows, err := db.Query(ctx, createSchedules, arg.Workflowrunid, arg.Triggertimes, arg.Input)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*WorkflowTriggerScheduledRef
	for rows.Next() {
		var i WorkflowTriggerScheduledRef
		if err := rows.Scan(
			&i.ID,
			&i.ParentId,
			&i.TriggerAt,
			&i.TickerId,
			&i.Input,
			&i.ChildIndex,
			&i.ChildKey,
			&i.ParentStepRunId,
			&i.ParentWorkflowRunId,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createStep = `-- name: CreateStep :one
INSERT INTO "Step" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "readableId",
    "tenantId",
    "jobId",
    "actionId",
    "timeout",
    "customUserData",
    "retries",
    "scheduleTimeout"
) VALUES (
    $1::uuid,
    coalesce($2::timestamp, CURRENT_TIMESTAMP),
    coalesce($3::timestamp, CURRENT_TIMESTAMP),
    $4::timestamp,
    $5::text,
    $6::uuid,
    $7::uuid,
    $8::text,
    $9::text,
    coalesce($10::jsonb, '{}'),
    coalesce($11::integer, 0),
    coalesce($12::text, '5m')
) RETURNING id, "createdAt", "updatedAt", "deletedAt", "readableId", "tenantId", "jobId", "actionId", timeout, "customUserData", retries, "scheduleTimeout"
`

type CreateStepParams struct {
	ID              pgtype.UUID      `json:"id"`
	CreatedAt       pgtype.Timestamp `json:"createdAt"`
	UpdatedAt       pgtype.Timestamp `json:"updatedAt"`
	Deletedat       pgtype.Timestamp `json:"deletedat"`
	Readableid      string           `json:"readableid"`
	Tenantid        pgtype.UUID      `json:"tenantid"`
	Jobid           pgtype.UUID      `json:"jobid"`
	Actionid        string           `json:"actionid"`
	Timeout         pgtype.Text      `json:"timeout"`
	CustomUserData  []byte           `json:"customUserData"`
	Retries         pgtype.Int4      `json:"retries"`
	ScheduleTimeout pgtype.Text      `json:"scheduleTimeout"`
}

func (q *Queries) CreateStep(ctx context.Context, db DBTX, arg CreateStepParams) (*Step, error) {
	row := db.QueryRow(ctx, createStep,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Deletedat,
		arg.Readableid,
		arg.Tenantid,
		arg.Jobid,
		arg.Actionid,
		arg.Timeout,
		arg.CustomUserData,
		arg.Retries,
		arg.ScheduleTimeout,
	)
	var i Step
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.ReadableId,
		&i.TenantId,
		&i.JobId,
		&i.ActionId,
		&i.Timeout,
		&i.CustomUserData,
		&i.Retries,
		&i.ScheduleTimeout,
	)
	return &i, err
}

const createStepRateLimit = `-- name: CreateStepRateLimit :one
INSERT INTO "StepRateLimit" (
    "units",
    "stepId",
    "rateLimitKey",
    "tenantId"
) VALUES (
    $1::integer,
    $2::uuid,
    $3::text,
    $4::uuid
) RETURNING units, "stepId", "rateLimitKey", "tenantId"
`

type CreateStepRateLimitParams struct {
	Units        int32       `json:"units"`
	Stepid       pgtype.UUID `json:"stepid"`
	Ratelimitkey string      `json:"ratelimitkey"`
	Tenantid     pgtype.UUID `json:"tenantid"`
}

func (q *Queries) CreateStepRateLimit(ctx context.Context, db DBTX, arg CreateStepRateLimitParams) (*StepRateLimit, error) {
	row := db.QueryRow(ctx, createStepRateLimit,
		arg.Units,
		arg.Stepid,
		arg.Ratelimitkey,
		arg.Tenantid,
	)
	var i StepRateLimit
	err := row.Scan(
		&i.Units,
		&i.StepId,
		&i.RateLimitKey,
		&i.TenantId,
	)
	return &i, err
}

const createWorkflow = `-- name: CreateWorkflow :one
INSERT INTO "Workflow" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "tenantId",
    "name",
    "description"
) VALUES (
    $1::uuid,
    coalesce($2::timestamp, CURRENT_TIMESTAMP),
    coalesce($3::timestamp, CURRENT_TIMESTAMP),
    $4::timestamp,
    $5::uuid,
    $6::text,
    $7::text
) RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", name, description
`

type CreateWorkflowParams struct {
	ID          pgtype.UUID      `json:"id"`
	CreatedAt   pgtype.Timestamp `json:"createdAt"`
	UpdatedAt   pgtype.Timestamp `json:"updatedAt"`
	Deletedat   pgtype.Timestamp `json:"deletedat"`
	Tenantid    pgtype.UUID      `json:"tenantid"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
}

func (q *Queries) CreateWorkflow(ctx context.Context, db DBTX, arg CreateWorkflowParams) (*Workflow, error) {
	row := db.QueryRow(ctx, createWorkflow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Deletedat,
		arg.Tenantid,
		arg.Name,
		arg.Description,
	)
	var i Workflow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.Name,
		&i.Description,
	)
	return &i, err
}

const createWorkflowConcurrency = `-- name: CreateWorkflowConcurrency :one
INSERT INTO "WorkflowConcurrency" (
    "id",
    "createdAt",
    "updatedAt",
    "workflowVersionId",
    "getConcurrencyGroupId",
    "maxRuns",
    "limitStrategy"
) VALUES (
    $1::uuid,
    coalesce($2::timestamp, CURRENT_TIMESTAMP),
    coalesce($3::timestamp, CURRENT_TIMESTAMP),
    $4::uuid,
    $5::uuid,
    coalesce($6::integer, 1),
    coalesce($7::"ConcurrencyLimitStrategy", 'CANCEL_IN_PROGRESS')
) RETURNING id, "createdAt", "updatedAt", "workflowVersionId", "getConcurrencyGroupId", "maxRuns", "limitStrategy"
`

type CreateWorkflowConcurrencyParams struct {
	ID                    pgtype.UUID                  `json:"id"`
	CreatedAt             pgtype.Timestamp             `json:"createdAt"`
	UpdatedAt             pgtype.Timestamp             `json:"updatedAt"`
	Workflowversionid     pgtype.UUID                  `json:"workflowversionid"`
	Getconcurrencygroupid pgtype.UUID                  `json:"getconcurrencygroupid"`
	MaxRuns               pgtype.Int4                  `json:"maxRuns"`
	LimitStrategy         NullConcurrencyLimitStrategy `json:"limitStrategy"`
}

func (q *Queries) CreateWorkflowConcurrency(ctx context.Context, db DBTX, arg CreateWorkflowConcurrencyParams) (*WorkflowConcurrency, error) {
	row := db.QueryRow(ctx, createWorkflowConcurrency,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Workflowversionid,
		arg.Getconcurrencygroupid,
		arg.MaxRuns,
		arg.LimitStrategy,
	)
	var i WorkflowConcurrency
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.WorkflowVersionId,
		&i.GetConcurrencyGroupId,
		&i.MaxRuns,
		&i.LimitStrategy,
	)
	return &i, err
}

const createWorkflowTriggerCronRef = `-- name: CreateWorkflowTriggerCronRef :one
INSERT INTO "WorkflowTriggerCronRef" (
    "parentId",
    "cron",
    "input"
) VALUES (
    $1::uuid,
    $2::text,
    $3::jsonb
) RETURNING "parentId", cron, "tickerId", input, enabled
`

type CreateWorkflowTriggerCronRefParams struct {
	Workflowtriggersid pgtype.UUID `json:"workflowtriggersid"`
	Crontrigger        string      `json:"crontrigger"`
	Input              []byte      `json:"input"`
}

func (q *Queries) CreateWorkflowTriggerCronRef(ctx context.Context, db DBTX, arg CreateWorkflowTriggerCronRefParams) (*WorkflowTriggerCronRef, error) {
	row := db.QueryRow(ctx, createWorkflowTriggerCronRef, arg.Workflowtriggersid, arg.Crontrigger, arg.Input)
	var i WorkflowTriggerCronRef
	err := row.Scan(
		&i.ParentId,
		&i.Cron,
		&i.TickerId,
		&i.Input,
		&i.Enabled,
	)
	return &i, err
}

const createWorkflowTriggerEventRef = `-- name: CreateWorkflowTriggerEventRef :one
INSERT INTO "WorkflowTriggerEventRef" (
    "parentId",
    "eventKey"
) VALUES (
    $1::uuid,
    $2::text
) RETURNING "parentId", "eventKey"
`

type CreateWorkflowTriggerEventRefParams struct {
	Workflowtriggersid pgtype.UUID `json:"workflowtriggersid"`
	Eventtrigger       string      `json:"eventtrigger"`
}

func (q *Queries) CreateWorkflowTriggerEventRef(ctx context.Context, db DBTX, arg CreateWorkflowTriggerEventRefParams) (*WorkflowTriggerEventRef, error) {
	row := db.QueryRow(ctx, createWorkflowTriggerEventRef, arg.Workflowtriggersid, arg.Eventtrigger)
	var i WorkflowTriggerEventRef
	err := row.Scan(&i.ParentId, &i.EventKey)
	return &i, err
}

const createWorkflowTriggerScheduledRef = `-- name: CreateWorkflowTriggerScheduledRef :one
INSERT INTO "WorkflowTriggerScheduledRef" (
    "id",
    "parentId",
    "triggerAt",
    "tickerId",
    "input"
) VALUES (
    gen_random_uuid(),
    $1::uuid,
    $2::timestamp,
    NULL, -- or provide a tickerId if applicable
    NULL -- or provide input if applicable
) RETURNING id, "parentId", "triggerAt", "tickerId", input, "childIndex", "childKey", "parentStepRunId", "parentWorkflowRunId"
`

type CreateWorkflowTriggerScheduledRefParams struct {
	Workflowversionid pgtype.UUID      `json:"workflowversionid"`
	Scheduledtrigger  pgtype.Timestamp `json:"scheduledtrigger"`
}

func (q *Queries) CreateWorkflowTriggerScheduledRef(ctx context.Context, db DBTX, arg CreateWorkflowTriggerScheduledRefParams) (*WorkflowTriggerScheduledRef, error) {
	row := db.QueryRow(ctx, createWorkflowTriggerScheduledRef, arg.Workflowversionid, arg.Scheduledtrigger)
	var i WorkflowTriggerScheduledRef
	err := row.Scan(
		&i.ID,
		&i.ParentId,
		&i.TriggerAt,
		&i.TickerId,
		&i.Input,
		&i.ChildIndex,
		&i.ChildKey,
		&i.ParentStepRunId,
		&i.ParentWorkflowRunId,
	)
	return &i, err
}

const createWorkflowTriggers = `-- name: CreateWorkflowTriggers :one
INSERT INTO "WorkflowTriggers" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "workflowVersionId",
    "tenantId"
) VALUES (
    $1::uuid,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    NULL,
    $2::uuid,
    $3::uuid
) RETURNING id, "createdAt", "updatedAt", "deletedAt", "workflowVersionId", "tenantId"
`

type CreateWorkflowTriggersParams struct {
	ID                pgtype.UUID `json:"id"`
	Workflowversionid pgtype.UUID `json:"workflowversionid"`
	Tenantid          pgtype.UUID `json:"tenantid"`
}

func (q *Queries) CreateWorkflowTriggers(ctx context.Context, db DBTX, arg CreateWorkflowTriggersParams) (*WorkflowTriggers, error) {
	row := db.QueryRow(ctx, createWorkflowTriggers, arg.ID, arg.Workflowversionid, arg.Tenantid)
	var i WorkflowTriggers
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.WorkflowVersionId,
		&i.TenantId,
	)
	return &i, err
}

const createWorkflowVersion = `-- name: CreateWorkflowVersion :one
INSERT INTO "WorkflowVersion" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "checksum",
    "version",
    "workflowId",
    "scheduleTimeout",
    "sticky",
    "kind",
    "defaultPriority"
) VALUES (
    $1::uuid,
    coalesce($2::timestamp, CURRENT_TIMESTAMP),
    coalesce($3::timestamp, CURRENT_TIMESTAMP),
    $4::timestamp,
    $5::text,
    $6::text,
    $7::uuid,
    coalesce($8::text, '5m'),
    $9::"StickyStrategy",
    coalesce($10::"WorkflowKind", 'DAG'),
    $11::integer
) RETURNING id, "createdAt", "updatedAt", "deletedAt", version, "order", "workflowId", checksum, "scheduleTimeout", "onFailureJobId", sticky, kind, "defaultPriority"
`

type CreateWorkflowVersionParams struct {
	ID              pgtype.UUID        `json:"id"`
	CreatedAt       pgtype.Timestamp   `json:"createdAt"`
	UpdatedAt       pgtype.Timestamp   `json:"updatedAt"`
	Deletedat       pgtype.Timestamp   `json:"deletedat"`
	Checksum        string             `json:"checksum"`
	Version         pgtype.Text        `json:"version"`
	Workflowid      pgtype.UUID        `json:"workflowid"`
	ScheduleTimeout pgtype.Text        `json:"scheduleTimeout"`
	Sticky          NullStickyStrategy `json:"sticky"`
	Kind            NullWorkflowKind   `json:"kind"`
	DefaultPriority pgtype.Int4        `json:"defaultPriority"`
}

func (q *Queries) CreateWorkflowVersion(ctx context.Context, db DBTX, arg CreateWorkflowVersionParams) (*WorkflowVersion, error) {
	row := db.QueryRow(ctx, createWorkflowVersion,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Deletedat,
		arg.Checksum,
		arg.Version,
		arg.Workflowid,
		arg.ScheduleTimeout,
		arg.Sticky,
		arg.Kind,
		arg.DefaultPriority,
	)
	var i WorkflowVersion
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Version,
		&i.Order,
		&i.WorkflowId,
		&i.Checksum,
		&i.ScheduleTimeout,
		&i.OnFailureJobId,
		&i.Sticky,
		&i.Kind,
		&i.DefaultPriority,
	)
	return &i, err
}

const getWorkflowByName = `-- name: GetWorkflowByName :one
SELECT
    id, "createdAt", "updatedAt", "deletedAt", "tenantId", name, description
FROM
    "Workflow" as workflows
WHERE
    workflows."tenantId" = $1::uuid AND
    workflows."name" = $2::text AND
    workflows."deletedAt" IS NULL
`

type GetWorkflowByNameParams struct {
	Tenantid pgtype.UUID `json:"tenantid"`
	Name     string      `json:"name"`
}

func (q *Queries) GetWorkflowByName(ctx context.Context, db DBTX, arg GetWorkflowByNameParams) (*Workflow, error) {
	row := db.QueryRow(ctx, getWorkflowByName, arg.Tenantid, arg.Name)
	var i Workflow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.Name,
		&i.Description,
	)
	return &i, err
}

const getWorkflowLatestVersion = `-- name: GetWorkflowLatestVersion :one
SELECT
    "id"
FROM
    "WorkflowVersion" as workflowVersions
WHERE
    workflowVersions."workflowId" = $1::uuid AND
    workflowVersions."deletedAt" IS NULL
ORDER BY
    workflowVersions."order" DESC
LIMIT 1
`

func (q *Queries) GetWorkflowLatestVersion(ctx context.Context, db DBTX, workflowid pgtype.UUID) (pgtype.UUID, error) {
	row := db.QueryRow(ctx, getWorkflowLatestVersion, workflowid)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const getWorkflowVersionForEngine = `-- name: GetWorkflowVersionForEngine :many
SELECT
    workflowversions.id, workflowversions."createdAt", workflowversions."updatedAt", workflowversions."deletedAt", workflowversions.version, workflowversions."order", workflowversions."workflowId", workflowversions.checksum, workflowversions."scheduleTimeout", workflowversions."onFailureJobId", workflowversions.sticky, workflowversions.kind, workflowversions."defaultPriority",
    w."name" as "workflowName",
    wc."limitStrategy" as "concurrencyLimitStrategy",
    wc."maxRuns" as "concurrencyMaxRuns"
FROM
    "WorkflowVersion" as workflowVersions
JOIN
    "Workflow" as w ON w."id" = workflowVersions."workflowId"
LEFT JOIN
    "WorkflowConcurrency" as wc ON wc."workflowVersionId" = workflowVersions."id"
WHERE
    workflowVersions."id" = ANY($1::uuid[]) AND
    w."tenantId" = $2::uuid AND
    w."deletedAt" IS NULL AND
    workflowVersions."deletedAt" IS NULL
`

type GetWorkflowVersionForEngineParams struct {
	Ids      []pgtype.UUID `json:"ids"`
	Tenantid pgtype.UUID   `json:"tenantid"`
}

type GetWorkflowVersionForEngineRow struct {
	WorkflowVersion          WorkflowVersion              `json:"workflow_version"`
	WorkflowName             string                       `json:"workflowName"`
	ConcurrencyLimitStrategy NullConcurrencyLimitStrategy `json:"concurrencyLimitStrategy"`
	ConcurrencyMaxRuns       pgtype.Int4                  `json:"concurrencyMaxRuns"`
}

func (q *Queries) GetWorkflowVersionForEngine(ctx context.Context, db DBTX, arg GetWorkflowVersionForEngineParams) ([]*GetWorkflowVersionForEngineRow, error) {
	rows, err := db.Query(ctx, getWorkflowVersionForEngine, arg.Ids, arg.Tenantid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetWorkflowVersionForEngineRow
	for rows.Next() {
		var i GetWorkflowVersionForEngineRow
		if err := rows.Scan(
			&i.WorkflowVersion.ID,
			&i.WorkflowVersion.CreatedAt,
			&i.WorkflowVersion.UpdatedAt,
			&i.WorkflowVersion.DeletedAt,
			&i.WorkflowVersion.Version,
			&i.WorkflowVersion.Order,
			&i.WorkflowVersion.WorkflowId,
			&i.WorkflowVersion.Checksum,
			&i.WorkflowVersion.ScheduleTimeout,
			&i.WorkflowVersion.OnFailureJobId,
			&i.WorkflowVersion.Sticky,
			&i.WorkflowVersion.Kind,
			&i.WorkflowVersion.DefaultPriority,
			&i.WorkflowName,
			&i.ConcurrencyLimitStrategy,
			&i.ConcurrencyMaxRuns,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkflowWorkerCount = `-- name: GetWorkflowWorkerCount :one
WITH UniqueWorkers AS (
    SELECT DISTINCT wss."id" AS slotId, w."id" AS wokerId, w."maxRuns" AS maxRuns
    FROM "Worker" w
    JOIN "_ActionToWorker" atw ON w."id" = atw."B"
    JOIN "Action" a ON atw."A" = a."id"
    JOIN "Step" s ON a."actionId" = s."actionId"
    JOIN "Job" j ON s."jobId" = j."id"
    JOIN "WorkflowVersion" workflowVersion ON j."workflowVersionId" = workflowVersion."id"
    JOIN "WorkerSemaphoreSlot" wss ON w."id" = wss."workerId"
    WHERE
        w."tenantId" = $1::uuid
        AND workflowVersion."deletedAt" IS NULL
        AND w."deletedAt" IS NULL
        AND w."dispatcherId" IS NOT NULL
        AND w."lastHeartbeatAt" > NOW() - INTERVAL '5 seconds'
        AND w."isActive" = true
        AND w."isPaused" = false
        AND workflowVersion."workflowId" = $2::uuid
        AND wss."stepRunId" IS NULL
),
    workers as ( Select SUM("maxRuns") as maxR from "Worker" where "id" in (select wokerId from UniqueWorkers)),
    slots as (
SELECT
    COUNT(uw.slotId) AS freeSlotCount
FROM UniqueWorkers uw)

SELECT
    maxR as totalCount,
    freeSlotCount as freeCount
FROM workers, slots
`

type GetWorkflowWorkerCountParams struct {
	Tenantid   pgtype.UUID `json:"tenantid"`
	Workflowid pgtype.UUID `json:"workflowid"`
}

type GetWorkflowWorkerCountRow struct {
	Totalcount int64 `json:"totalcount"`
	Freecount  int64 `json:"freecount"`
}

func (q *Queries) GetWorkflowWorkerCount(ctx context.Context, db DBTX, arg GetWorkflowWorkerCountParams) (*GetWorkflowWorkerCountRow, error) {
	row := db.QueryRow(ctx, getWorkflowWorkerCount, arg.Tenantid, arg.Workflowid)
	var i GetWorkflowWorkerCountRow
	err := row.Scan(&i.Totalcount, &i.Freecount)
	return &i, err
}

const linkOnFailureJob = `-- name: LinkOnFailureJob :one
UPDATE "WorkflowVersion"
SET "onFailureJobId" = $1::uuid
WHERE "id" = $2::uuid
RETURNING id, "createdAt", "updatedAt", "deletedAt", version, "order", "workflowId", checksum, "scheduleTimeout", "onFailureJobId", sticky, kind, "defaultPriority"
`

type LinkOnFailureJobParams struct {
	Jobid             pgtype.UUID `json:"jobid"`
	Workflowversionid pgtype.UUID `json:"workflowversionid"`
}

func (q *Queries) LinkOnFailureJob(ctx context.Context, db DBTX, arg LinkOnFailureJobParams) (*WorkflowVersion, error) {
	row := db.QueryRow(ctx, linkOnFailureJob, arg.Jobid, arg.Workflowversionid)
	var i WorkflowVersion
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Version,
		&i.Order,
		&i.WorkflowId,
		&i.Checksum,
		&i.ScheduleTimeout,
		&i.OnFailureJobId,
		&i.Sticky,
		&i.Kind,
		&i.DefaultPriority,
	)
	return &i, err
}

const listWorkflows = `-- name: ListWorkflows :many
SELECT
    workflows.id, workflows."createdAt", workflows."updatedAt", workflows."deletedAt", workflows."tenantId", workflows.name, workflows.description
FROM
    "Workflow" as workflows
WHERE
    workflows."tenantId" = $1 AND
    workflows."deletedAt" IS NULL
ORDER BY
    case when $2 = 'createdAt ASC' THEN workflows."createdAt" END ASC ,
    case when $2 = 'createdAt DESC' then workflows."createdAt" END DESC
OFFSET
    COALESCE($3, 0)
LIMIT
    COALESCE($4, 50)
`

type ListWorkflowsParams struct {
	TenantId pgtype.UUID `json:"tenantId"`
	Orderby  interface{} `json:"orderby"`
	Offset   interface{} `json:"offset"`
	Limit    interface{} `json:"limit"`
}

type ListWorkflowsRow struct {
	Workflow Workflow `json:"workflow"`
}

func (q *Queries) ListWorkflows(ctx context.Context, db DBTX, arg ListWorkflowsParams) ([]*ListWorkflowsRow, error) {
	rows, err := db.Query(ctx, listWorkflows,
		arg.TenantId,
		arg.Orderby,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListWorkflowsRow
	for rows.Next() {
		var i ListWorkflowsRow
		if err := rows.Scan(
			&i.Workflow.ID,
			&i.Workflow.CreatedAt,
			&i.Workflow.UpdatedAt,
			&i.Workflow.DeletedAt,
			&i.Workflow.TenantId,
			&i.Workflow.Name,
			&i.Workflow.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listWorkflowsForEvent = `-- name: ListWorkflowsForEvent :many
SELECT DISTINCT ON ("WorkflowVersion"."workflowId") "WorkflowVersion".id
FROM "WorkflowVersion"
LEFT JOIN "Workflow" AS j1 ON j1.id = "WorkflowVersion"."workflowId"
LEFT JOIN "WorkflowTriggers" AS j2 ON j2."workflowVersionId" = "WorkflowVersion"."id"
WHERE
    (j1."tenantId"::uuid = $1 AND j1.id IS NOT NULL)
    AND j1."deletedAt" IS NULL
    AND "WorkflowVersion"."deletedAt" IS NULL
    AND
    (j2.id IN (
        SELECT t3."parentId"
        FROM "WorkflowTriggerEventRef" AS t3
        WHERE t3."eventKey" = $2 AND t3."parentId" IS NOT NULL
    ) AND j2.id IS NOT NULL)
    AND "WorkflowVersion".id = (
        -- confirm that the workflow version is the latest
        SELECT wv2.id
        FROM "WorkflowVersion" wv2
        WHERE wv2."workflowId" = "WorkflowVersion"."workflowId"
        ORDER BY wv2."order" DESC
        LIMIT 1
    )
ORDER BY "WorkflowVersion"."workflowId", "WorkflowVersion"."order" DESC
`

type ListWorkflowsForEventParams struct {
	Tenantid pgtype.UUID `json:"tenantid"`
	Eventkey string      `json:"eventkey"`
}

func (q *Queries) ListWorkflowsForEvent(ctx context.Context, db DBTX, arg ListWorkflowsForEventParams) ([]pgtype.UUID, error) {
	rows, err := db.Query(ctx, listWorkflowsForEvent, arg.Tenantid, arg.Eventkey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.UUID
	for rows.Next() {
		var id pgtype.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listWorkflowsLatestRuns = `-- name: ListWorkflowsLatestRuns :many
SELECT
    DISTINCT ON (workflow."id") runs."createdAt", runs."updatedAt", runs."deletedAt", runs."tenantId", runs."workflowVersionId", runs.status, runs.error, runs."startedAt", runs."finishedAt", runs."concurrencyGroupId", runs."displayName", runs.id, runs."childIndex", runs."childKey", runs."parentId", runs."parentStepRunId", runs."additionalMetadata", runs.duration, runs.priority, workflow."id" as "workflowId"
FROM
    "WorkflowRun" as runs
LEFT JOIN
    "WorkflowVersion" as workflowVersion ON runs."workflowVersionId" = workflowVersion."id"
LEFT JOIN
    "Workflow" as workflow ON workflowVersion."workflowId" = workflow."id"
WHERE
    runs."tenantId" = $1 AND
    runs."deletedAt" IS NULL AND
    workflow."deletedAt" IS NULL AND
    workflowVersion."deletedAt" IS NULL AND
    (
        $2::text IS NULL OR
        workflow."id" IN (
            SELECT
                DISTINCT ON(t1."workflowId") t1."workflowId"
            FROM
                "WorkflowVersion" AS t1
                LEFT JOIN "WorkflowTriggers" AS j2 ON j2."workflowVersionId" = t1."id"
            WHERE
                (
                    j2."id" IN (
                        SELECT
                            t3."parentId"
                        FROM
                            "public"."WorkflowTriggerEventRef" AS t3
                        WHERE
                            t3."eventKey" = $2::text
                            AND t3."parentId" IS NOT NULL
                    )
                    AND j2."id" IS NOT NULL
                    AND t1."workflowId" IS NOT NULL
                )
            ORDER BY
                t1."workflowId" DESC, t1."order" DESC
        )
    )
ORDER BY
    workflow."id" DESC, runs."createdAt" DESC
`

type ListWorkflowsLatestRunsParams struct {
	TenantId pgtype.UUID `json:"tenantId"`
	EventKey pgtype.Text `json:"eventKey"`
}

type ListWorkflowsLatestRunsRow struct {
	WorkflowRun WorkflowRun `json:"workflow_run"`
	WorkflowId  pgtype.UUID `json:"workflowId"`
}

func (q *Queries) ListWorkflowsLatestRuns(ctx context.Context, db DBTX, arg ListWorkflowsLatestRunsParams) ([]*ListWorkflowsLatestRunsRow, error) {
	rows, err := db.Query(ctx, listWorkflowsLatestRuns, arg.TenantId, arg.EventKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListWorkflowsLatestRunsRow
	for rows.Next() {
		var i ListWorkflowsLatestRunsRow
		if err := rows.Scan(
			&i.WorkflowRun.CreatedAt,
			&i.WorkflowRun.UpdatedAt,
			&i.WorkflowRun.DeletedAt,
			&i.WorkflowRun.TenantId,
			&i.WorkflowRun.WorkflowVersionId,
			&i.WorkflowRun.Status,
			&i.WorkflowRun.Error,
			&i.WorkflowRun.StartedAt,
			&i.WorkflowRun.FinishedAt,
			&i.WorkflowRun.ConcurrencyGroupId,
			&i.WorkflowRun.DisplayName,
			&i.WorkflowRun.ID,
			&i.WorkflowRun.ChildIndex,
			&i.WorkflowRun.ChildKey,
			&i.WorkflowRun.ParentId,
			&i.WorkflowRun.ParentStepRunId,
			&i.WorkflowRun.AdditionalMetadata,
			&i.WorkflowRun.Duration,
			&i.WorkflowRun.Priority,
			&i.WorkflowId,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const softDeleteWorkflow = `-- name: SoftDeleteWorkflow :one
WITH versions AS (
    UPDATE "WorkflowVersion"
    SET "deletedAt" = CURRENT_TIMESTAMP
    WHERE "workflowId" = $1::uuid
)
UPDATE "Workflow"
SET
    -- set name to the current name plus a random suffix to avoid conflicts
    "name" = "name" || '-' || gen_random_uuid(),
    "deletedAt" = CURRENT_TIMESTAMP
WHERE "id" = $1::uuid
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", name, description
`

func (q *Queries) SoftDeleteWorkflow(ctx context.Context, db DBTX, id pgtype.UUID) (*Workflow, error) {
	row := db.QueryRow(ctx, softDeleteWorkflow, id)
	var i Workflow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.Name,
		&i.Description,
	)
	return &i, err
}

const upsertAction = `-- name: UpsertAction :one
INSERT INTO "Action" (
    "id",
    "actionId",
    "tenantId"
)
VALUES (
    gen_random_uuid(),
    LOWER($1::text),
    $2::uuid
)
ON CONFLICT ("tenantId", "actionId") DO UPDATE
SET
    "tenantId" = EXCLUDED."tenantId"
WHERE
    "Action"."tenantId" = $2 AND "Action"."actionId" = LOWER($1::text)
RETURNING description, "tenantId", "actionId", id
`

type UpsertActionParams struct {
	Action   string      `json:"action"`
	Tenantid pgtype.UUID `json:"tenantid"`
}

func (q *Queries) UpsertAction(ctx context.Context, db DBTX, arg UpsertActionParams) (*Action, error) {
	row := db.QueryRow(ctx, upsertAction, arg.Action, arg.Tenantid)
	var i Action
	err := row.Scan(
		&i.Description,
		&i.TenantId,
		&i.ActionId,
		&i.ID,
	)
	return &i, err
}

const upsertWorkflowTag = `-- name: UpsertWorkflowTag :exec
INSERT INTO "WorkflowTag" (
    "id",
    "tenantId",
    "name",
    "color"
)
VALUES (
    COALESCE($1::uuid, gen_random_uuid()),
    $2::uuid,
    $3::text,
    COALESCE($4::text, '#93C5FD')
)
ON CONFLICT ("tenantId", "name") DO UPDATE
SET
    "color" = COALESCE(EXCLUDED."color", "WorkflowTag"."color")
WHERE
    "WorkflowTag"."tenantId" = $2 AND "WorkflowTag"."name" = $3
`

type UpsertWorkflowTagParams struct {
	ID       pgtype.UUID `json:"id"`
	Tenantid pgtype.UUID `json:"tenantid"`
	Tagname  string      `json:"tagname"`
	TagColor pgtype.Text `json:"tagColor"`
}

func (q *Queries) UpsertWorkflowTag(ctx context.Context, db DBTX, arg UpsertWorkflowTagParams) error {
	_, err := db.Exec(ctx, upsertWorkflowTag,
		arg.ID,
		arg.Tenantid,
		arg.Tagname,
		arg.TagColor,
	)
	return err
}
