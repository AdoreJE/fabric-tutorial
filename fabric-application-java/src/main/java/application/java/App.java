/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

// Running TestApp:
// gradle runApp

package application.java;

import java.nio.file.Path;
import java.nio.file.Paths;

import org.hyperledger.fabric.gateway.*;


public class App {

	private static final String CHANNEL_NAME = System.getenv().getOrDefault("CHANNEL_NAME", "mychannel");
	private static final String CHAINCODE_NAME = System.getenv().getOrDefault("CHAINCODE_NAME", "basic");

	static {
		System.setProperty("org.hyperledger.fabric.sdk.service_discovery.as_localhost", "false");
	}

	// helper function for getting connected to the gateway
	public static Gateway connect() throws Exception{
		// Load a file system based wallet for managing identities.
		Path walletPath = Paths.get("wallet");
		Wallet wallet = Wallets.newFileSystemWallet(walletPath);
		// load a CCP
		Path networkConfigPath = Paths.get("connection-org1.yaml");

		Gateway.Builder builder = Gateway.createBuilder();
		builder.identity(wallet, "att3").networkConfig(networkConfigPath).discovery(true);

		return builder.connect();
	}

	public static void main(String[] args) throws Exception {
		// enrolls the admin and registers the user

		// connect to the network and invoke the smart contract
		try (Gateway gateway = connect()) {

			// get the network and contract
			Network network = gateway.getNetwork(CHANNEL_NAME);
			Contract contract = network.getContract("ping");

			byte[] result;

			//System.out.println("Submit Transaction: InitLedger creates the initial set of assets on the ledger.");
			//contract.submitTransaction("InitLedger");

			System.out.println("\n");
			result = contract.evaluateTransaction("ClientInfo");
			System.out.println("Evaluate Transaction: GetAllAssets, result: " + new String(result));
		}
		catch(Exception e){
			System.err.println(e);
			System.exit(1);
		}
	}
}