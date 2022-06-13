import {COMPONENT_TEST_TIMEOUT, MONGO_DSN, SERVICE_HOST} from '../config'
import {sleep} from '../common/helpers'
import {AfterAll, Before, BeforeAll, setDefaultTimeout} from 'cucumber'
import * as Amqp from 'amqp-ts'
import * as AxiosLogger from 'axios-logger'
import axios from 'axios'
import {requestLogger} from '../common/logging'

setDefaultTimeout(Math.max(+COMPONENT_TEST_TIMEOUT, 10) * 1000)

axios.interceptors.request.use(AxiosLogger.requestLogger)
axios.interceptors.response.use(AxiosLogger.responseLogger)
const axiosInstance = axios.create()

export const testQueues: Amqp.Queue[] = []

async function pollService(address: string): Promise<void> {
  console.log(
    `Waiting up to ${COMPONENT_TEST_TIMEOUT} seconds service ${address} to be up and ready \n\n`,
  )

  for (let i = 0; i < COMPONENT_TEST_TIMEOUT; i++) {
    try {
      const response = await axiosInstance.get(address)
      console.log(
        'Service is up and responded: ' + response.status + '\n' + JSON.stringify(response.data),
      )

      return
    } catch (error) {
      console.error('Could not connect to service ' + error)
      await sleep(1000)
    }
  }

  throw new Error('Service did not start up within 4 seconds')
}

BeforeAll(async () => {
  await pollService(SERVICE_HOST + '/checks/readiness')
  console.log(' ==> Trying to connect to MongoDB @ ' + MONGO_DSN + '\n')
  try {
    // We can connect to MongoDB here
    console.log(' Successfully connected to MongoDB \n\n')
  } catch (err) {
    console.log('!!!! ERROR could not connect to MongoDB!!! \n\n\n\n')
    throw err
  }
})

Before(async function () {
  requestLogger.setWorld(this)
})

AfterAll(async () => {
  // Close connections
})
