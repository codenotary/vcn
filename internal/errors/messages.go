/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */
package errors

const AuthRequired = `authentication required, please login`

const UnconfirmedEmail = `your email address was not confirmed

Please confirm it by clicking on the link we sent to %s.
If you did not receive the email, please go to %s and click on the link "Resend email"`

const NoRemainingSignOps = `signature quota exceeded

Unfortunately, you have used all your signatures for this month.
If you would like to increase the number of monthly signatures, 
please email us at support@codenotary.io with your request.`

// BlockchainPermission refers to https://github.com/vchain-us/vcn/wiki/Errors#blockchain-permission-403
const BlockchainPermission = `could not write to blockchain

You have not permission to write onto the blockchain (yet).
Please try again later. If the problem persists contact our support.`

const BlockchainTimeout = `writing to blockchain timed out

Please try again later. If the problem persists contact our support.`

const BlockchainCannotConnect = `cannot connect to blockchain

Please try again later. If the problem persists contact our support.`

const BlockchainContractErr = `cannot instantiate contract

Please try again later. If the problem persists contact our support.`

const SignFailed = `method <sign> failed`
