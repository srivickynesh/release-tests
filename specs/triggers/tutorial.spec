PIPELINES-06
# Verify triggers tutorial

Pre condition:
  * Validate Operator should be installed

## Run pipelines tutorials: PIPELINES-06-TC01
Tags: e2e, integration, non-admin, pipelines, tutorial
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

This scenario tests the pipeline tutorials (https://github.com/openshift/pipelines-tutorial) pipelines related resources 

Steps:
  * Create remote
      |S.NO|resource_dir                                                                                                                             |
      |----|-----------------------------------------------------------------------------------------------------------------------------------------|
      |1   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/01_apply_manifest_task.yaml             |
      |2   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/02_update_deployment_task.yaml          |
      |3   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/03_persistent_volume_claim.yaml         |
      |4   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/04_pipeline.yaml                        |
      |5   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/02_pipelinerun/01_build_deploy_api_pipelinerun.yaml |
      |6   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/02_pipelinerun/02_build_deploy_ui_pipelinerun.yaml|
   * Verify pipelinerun
      |S.NO|pipeline_run_name           |status    |
      |----|----------------------------|----------|
      |1   |build-deploy-api-pipelinerun|successful|
      |2   |build-deploy-ui-pipelinerun |successful|
   * Get route url of the route "pipelines-vote-ui"
   * Wait for "pipelines-vote-api" deployment
   * Wait for "pipelines-vote-ui" deployment
   * Validate that route URL contains "Cat 🐺 vs Dog 🐶"

## Run pipelines tutorial using triggers: PIPELINES-06-TC02
Tags: e2e, integration, triggers, non-admin, tutorial, sanity
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario tests the pipeline tutorials (https://github.com/openshift/pipelines-tutorial) triggers related resources 

Steps:
  * Create remote
      |S.NO|resource_dir                                                                                                                             |
      |----|-----------------------------------------------------------------------------------------------------------------------------------------|
      |1   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/01_apply_manifest_task.yaml             |
      |2   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/02_update_deployment_task.yaml          |
      |3   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/03_persistent_volume_claim.yaml         |
      |4   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/04_pipeline.yaml                        |
      |5   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/03_triggers/01_binding.yaml                         |
      |6   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/03_triggers/02_template.yaml                        |
      |7   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/03_triggers/03_trigger.yaml                         |
      |8   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/03_triggers/04_event_listener.yaml                  |
   * Expose Event listener "vote-app"
   * Mock post event to "github" interceptor with event-type "push", payload "testdata/push-vote-api.json", with TLS "false"
   * Assert eventlistener response
   * "1" pipelinerun(s) should be present within "15" seconds
   * Verify the latest pipelinerun for "successful" state
   * Mock post event to "github" interceptor with event-type "push", payload "testdata/push-vote-ui.json", with TLS "false"
   * Assert eventlistener response
   * "2" pipelinerun(s) should be present within "15" seconds
   * Verify the latest pipelinerun for "successful" state
   * Get route url of the route "pipelines-vote-ui"
   * Wait for "pipelines-vote-ui" deployment
   * Validate that route URL contains "Cat 🐺 vs Dog 🐶"
