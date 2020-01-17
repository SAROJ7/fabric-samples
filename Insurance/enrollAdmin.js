

'use strict';
const fabricCAServices = require('fabric-ca-client');
const { FileSystemWallet, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(_dirname,'..','connection-org1.json');
const ccpJSON = fs.readFileSync()
