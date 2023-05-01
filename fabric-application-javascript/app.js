'use strict';

const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const path = require('path');
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('./CAUtil.js');
const { buildCCPOrg1, buildWallet, prettyJSONString } = require('./AppUtil.js');

async function main() {
    const ccp = buildCCPOrg1();
    const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');

    const walletPath = path.join(__dirname, 'wallet');
    const wallet = await buildWallet(Wallets, walletPath);

    await enrollAdmin(caClient, wallet, "Org1MSP");

    await registerAndEnrollUser(caClient, wallet, "Org1MSP", "appUser", 'org1.department1');

    const gateway = new Gateway();
    try {
        await gateway.connect(ccp, {
            wallet,
            identity: "appUser",
            discovery: { enabled: true, asLocalhost: false }
        });

        const network = await gateway.getNetwork("mychannel");

        const contract = network.getContract("ledger");

        //await contract.submitTransaction('Join', "user1", "name1", "21", "M");
        //await contract.submitTransaction('Join', "user2", "name2", "22", "M");
        //await contract.submitTransaction('Join', "user3", "name3", "23", "F");
        //await contract.submitTransaction('Join', "user4", "name4", "24", "F");

        //await contract.submitTransaction('MakeOrder', "order1", "user1", "item1");
        //await contract.submitTransaction('MakeOrder', "order2", "user4", "item2");

        let queryString = {
            selector : {
                "docType" : "order",
                "item" : {
                    "$eq": "item1"
                }
            }
        }
        let userList = await contract.evaluateTransaction('QueryStringWithItem', JSON.stringify(queryString));
        console.log(userList.toString())
    } finally {
        gateway.disconnect();
    }
}

main();