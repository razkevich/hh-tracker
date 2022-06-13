import {binding} from 'cucumber-tsflow/dist'
import {TestContext} from '../shared/common/test-context'

@binding([TestContext])
export class MessagesSteps {
  constructor(protected testContext: TestContext) {}
}
