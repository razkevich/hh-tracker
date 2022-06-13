import {AxiosResponse} from 'axios'
import {throwIfNullOrUndefined} from './helpers'
import {Queue, Connection} from 'amqp-ts'
import {testQueues} from '../hooks/setup'

export class TestContext {
  private _latestResponse: AxiosResponse | undefined
  private _responsesStack: AxiosResponse[] = []
  private _messagesConnection: Connection | undefined
  private _latestMessage: any
  private _messagesStack: any[] = []
  private _savedMessagesCount = 0

  get messagesConnection(): Connection | undefined {
    return this._messagesConnection
  }

  set messagesConnection(value: Connection | undefined) {
    this._messagesConnection = value
  }

  public ensureLatestResponse(): AxiosResponse {
    const latestResponse = this._latestResponse
    throwIfNullOrUndefined(
      latestResponse,
      'Test context last response not found. This is likely a bug.',
    )
    return latestResponse!
  }

  public ensurePreviousResponse(): AxiosResponse {
    const previousResponse = this._responsesStack[this._responsesStack.length - 2]
    throwIfNullOrUndefined(
      previousResponse,
      'Test context previous response not found. This is likely a bug.',
    )
    return previousResponse!
  }

  public ensureLatestMessage(): any {
    const latestMessage = this._latestMessage
    throwIfNullOrUndefined(
      latestMessage,
      'Test context latest response not found. This is likely a bug.',
    )
    return latestMessage!
  }

  public addMessage(message: any) {
    this._latestMessage = message
    this._messagesStack.push(message)
  }

  public saveCurrentReceivedMessagesAmount() {
    this._savedMessagesCount = this._messagesStack.length
  }

  get savedMessagesCount(): number {
    return this._savedMessagesCount
  }

  public getCurrentMessagesAmount(): number {
    return this._messagesStack.length
  }

  public addResponse(response: AxiosResponse): void {
    this._latestResponse = response
    this._responsesStack.push(response)
  }

  public getTestQueueByName(queueName: string): Queue {
    const foundQueue = testQueues.find((queue) => {
      return queue.name === queueName
    })
    throwIfNullOrUndefined(foundQueue, `Test queue with name ${queueName} was not found`)
    return foundQueue
  }
}
