/* eslint-disable @typescript-eslint/camelcase */
import {binding} from 'cucumber-tsflow'
import {TestContext} from '../shared/common/test-context'

@binding([TestContext])
export class MockSteps {
  constructor(protected testContext: TestContext) {}
}
