import {TestContext} from '../shared/common/test-context'
import {binding, given} from 'cucumber-tsflow/dist'
import {resetMocks} from '../shared/mocks/api'
import {sleep} from '../shared/common/helpers'

@binding([TestContext])
export class CommonSteps {
  constructor(protected testContext: TestContext) {}

  @given('I reset DB and mocks')
  public async resetDb() {
    // todo clear mongo
    await resetMocks()
  }

  @given('I reset mocks')
  public async resetMocks() {
    await resetMocks()
  }

  @given('I sleep for {int} milliseconds')
  public async sleep(ms: number) {
    await sleep(ms)
  }
}
