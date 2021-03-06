const ethers = require("ethers");

const ValidatorActionType = {
    Add: 0,
    Remove: 1
};

const Vote = {
    Yes: 0,
    No: 1,
}

const VoteStatus = {
    Inactive: 0,
    Active: 1,
    Finalized: 2
}

const ThresholdType = {
    Validator: 0,
    Deposit: 1,
}

// TODO format differently
const dummyData = {
    type: "erc20",
    value: 100,
    to: "0x1",
    from: "0x2"
}

const CreateDepositData = (data = dummyData, depositId = 0, originChain = 0) => {
    const hash = ethers.utils.keccak256(ethers.utils.toUtf8Bytes(data));
    return [
        hash,
        depositId,
        originChain,
    ];
}

module.exports = {
    ValidatorActionType,
    Vote,
    VoteStatus,
    CreateDepositData,
    ThresholdType,
}