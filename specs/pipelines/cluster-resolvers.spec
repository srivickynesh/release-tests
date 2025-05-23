PIPELINES-23
# Cluster resolvers spec

Pre condition:
  * Create project "releasetest-tasks"
  * Apply
    |S.NO|resource_dir                                                |
    |----|------------------------------------------------------------|
    |1   |testdata/resolvers/tasks/resolver-task.yaml                 |
    |2   |testdata/resolvers/tasks/resolver-task2.yaml                |
  * Create project "releasetest-pipelines"
  * Apply
    |S.NO|resource_dir                                                |
    |----|------------------------------------------------------------|
    |1   |testdata/resolvers/pipelines/resolver-pipeline.yaml         |

## Checking the functionality of cluster resolvers#1: PIPELINES-23-TC01
Tags: e2e, sanity
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps:
    * Switch to autogenerated namespace
    * Create
      |S.NO|resource_dir                                                              |
      |----|--------------------------------------------------------------------------|
      |1   |testdata/resolvers/pipelineruns/resolver-pipelinerun.yaml                 |
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |
      |----|-----------------------------------|------------|
      |1   |resolver-pipelinerun               |successful  |
            
## Checking the functionality of cluster resolvers#2: PIPELINES-23-TC02
Tags: e2e
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps: 
    * Switch to autogenerated namespace
    * Create
      |S.NO|resource_dir                                                              |
      |----|--------------------------------------------------------------------------|
      |1   |testdata/resolvers/pipelines/resolver-pipeline-same-ns.yaml               |
      |2   |testdata/resolvers/pipelineruns/resolver-pipelinerun-same-ns.yaml         |
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |
      |----|-----------------------------------|------------|
      |1   |resolver-pipelinerun-same-ns       |successful  |
  
Teardown:
  * Delete project "releasetest-tasks"
  * Delete project "releasetest-pipelines"