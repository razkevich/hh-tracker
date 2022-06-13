import {TestContext} from '../shared/common/test-context'
import {binding} from 'cucumber-tsflow/dist'

@binding([TestContext])
export class CommonSteps {
  constructor(protected testContext: TestContext) {}
}
