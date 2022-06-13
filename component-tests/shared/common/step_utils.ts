import {expect} from 'chai'
import {TestContext} from './test-context'
import {binding} from 'cucumber-tsflow/dist'

@binding([TestContext])
export class StepUtils {
  constructor(protected testContext: TestContext) {}

  public assertResponse(code: number): void {
    expect(this.testContext.ensureLatestResponse().status).to.equal(
      code,
      'Unexpected status code in the latest response; response data: ' +
        JSON.stringify(this.testContext.ensureLatestResponse().data),
    )
  }
}
