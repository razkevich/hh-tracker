import {TestContext} from '../shared/common/test-context'
import {binding, when} from 'cucumber-tsflow/dist'
import {use} from 'chai'
import chaiDeepEqualInAnyOrder from 'deep-equal-in-any-order'
import {StepUtils} from '../shared/common/step_utils'
import {api} from '../shared/common/api'
import {LOG_PATH} from '../shared/logs/log-config'

use(chaiDeepEqualInAnyOrder)

@binding([TestContext])
export class LogsSteps {
  constructor(protected testContext: TestContext) {}

  private _stepUtils = new StepUtils(this.testContext)

  @when('I read logs for resource type {string} and store {string}', '', 60000)
  public async readLogs(resourceType: string, store: string) {
    this.testContext.addResponse(await api.get(LOG_PATH, store))
  }
}
