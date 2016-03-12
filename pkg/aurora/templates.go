package aurora

var DefaultJobJSON = `
{
    "cluster": "{{.cluster}}",
    "contact": "{{.contact}}",
    "container": {
        "docker": {
            "image": "{{.compose.image}}",
            "parameters": [
                {
                    "name": "l",
                    "value": "cluster.scheduler=aurora"
                }{{range .docker_params}}
                ,{
                    "name": "{{.Key}}",
                    "value": "{{.Value}}"
                }{{end}}]
        }
    },
    "enable_hooks": false,
    "environment": "{{.environment}}",
    "instances": 1,
    "max_task_failures": 1,
    "name": "{{.name}}",
    "priority": 0,
    "production": {{.is_production}},
    "role": "{{.role}}",
    "service": {{.is_service}},
    "task": {
        "constraints": [],
        "finalization_wait": 30,
        "max_concurrency": 0,
        "max_failures": 1,
        "name": "{{.name}}",
        "processes": [],
        "resources": {
            "cpu": 0.4,
            "disk": 67108864,
            "ram": 33554432
        }
    },
    "update_config": {
        "batch_size": 1,
        "max_per_shard_failures": 0,
        "max_total_failures": 0,
        "rollback_on_failure": true,
        "wait_for_batch_completion": false,
        "watch_secs": 45
    }
}
`
