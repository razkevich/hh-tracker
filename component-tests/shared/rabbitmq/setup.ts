import {rabbitClientWiremock} from './client'
import {AfterAll, Before, BeforeAll} from 'cucumber'
import {testQueues} from '../hooks/setup'

Before(async function () {
  if (testQueues !== undefined) {
    for (const queue of testQueues) {
      await rabbitClientWiremock.deleteQueue(queue)
    }
    testQueues.length = 0
  }
})

BeforeAll(async () => {
  await rabbitClientWiremock.waitForConnection()
})

AfterAll(async () => {
  await rabbitClientWiremock.closeConnection()
})
