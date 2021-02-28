


const Footsites = require('./normal/footsites');
const FootsitesWise = require('./wise/footsites');
const FootsitesProxyKiller = require('./proxy-killer/footsites.beast');
const FootsitesProxyAdvanced = require('./advance/footsites');
// const FOOTLOCKERCA = require('./FOOTLOCKERCA');

// const FootsitesCA = require('./normal/footlockerCa');
// const FootsitesProxyKillerCA = require('./proxy-killer/footlockerCa.beast');
// const FootsitesProxyAdvancedCA = require('./advance/footlockerCa');

// const FootsitesBeast = require('./footsites.beast');
// const FOOTLOCKERCABeast = require('./FOOTLOCKERCA.beast');

const runFootsites = async (task, mainWindow) => {

    let task1 = {
        id: task.id,
        site: task.store.longCode.toLowerCase(),
        pid: task.keywords,
        size: task.size,
        monitorDelay: task.advanced.delays.monitor.min != '' ? task.advanced.delays.monitor.min : 3000,
        navigationDelay: 0,
        restockDelay: task.advanced.delays.monitor.min != '' ? task.advanced.delays.monitor.min : 3000,
        proxy: task.proxy ? task.proxy : null,
        profile: {
            email: task.profile.paymentInfo.email,
            shippingAddress: {
                setAsDefaultBilling: false,
                setAsDefaultShipping: false,
                firstName: task.profile.shippingInfo.firstName,
                lastName: task.profile.shippingInfo.lastName,
                email: false,
                phone: task.profile.shippingInfo.phoneNumber,
                country: {
                    isocode: task.store.longCode.toUpperCase() === 'FOOTLOCKERCA' ? "CA" : "US",
                    name: task.store.longCode.toUpperCase() === 'FOOTLOCKERCA' ? "Canada" : "United States",
                },
                id: null,
                setAsBilling: true,
                region: {
                    countryIso: task.store.longCode.toUpperCase() === 'FOOTLOCKERCA' ? "CA" : "US",
                    isocode: task.store.longCode.toUpperCase() === 'FOOTLOCKERCA' ? `CA-${task.profile.shippingInfo.state}` : `US-${task.profile.shippingInfo.state}`,
                    isocodeShort: `${task.profile.shippingInfo.state}`,
                    name: `${abbrState(task.profile.shippingInfo.state, 'name')}`,
                },
                type: "default",
                LoqateSearch: "",
                line1: `${task.profile.shippingInfo.address}`,
                postalCode: `${task.profile.shippingInfo.zipcode}`,
                town: `${task.profile.shippingInfo.city}`,
                regionFPO: null,
                shippingAddress: true,
                recordType: "S",
            },
            billingAddress: {
                setAsDefaultBilling: false,
                setAsDefaultShipping: false,
                firstName: task.profile.billingInfo.firstName,
                lastName: task.profile.billingInfo.lastName,
                email: false,
                phone: task.profile.billingInfo.phoneNumber,
                country: { isocode: task.store.longCode.toUpperCase() === 'FOOTLOCKERCA' ? "CA" : "US", name: "United States" },
                id: null,
                setAsBilling: false,
                region: {
                    countryIso: task.store.longCode.toUpperCase() === 'FOOTLOCKERCA' ? "CA" : "US",
                    isocode: task.store.longCode.toUpperCase() === 'FOOTLOCKERCA' ? `CA-${task.profile.shippingInfo.state}` : `US-${task.profile.shippingInfo.state}`,
                    isocodeShort: `${task.profile.shippingInfo.state}`,
                    name: `${abbrState(task.profile.shippingInfo.state, 'name')}`,
                },
                type: "default",
                LoqateSearch: "",
                line1: task.profile.billingInfo.address,
                postalCode: task.profile.billingInfo.zipcode,
                town: task.profile.billingInfo.city,
                regionFPO: null,
                shippingAddress: true,
                recordType: "S",
            },
            paymentInfo: {
                cardNumber: task.profile.paymentInfo.cardNumber.match(/.{1,4}/g).join(" "),
                cardMonth: task.profile.paymentInfo.expiryMonth,
                cardYear: task.profile.paymentInfo.expiryYear,
                cardCvv: task.profile.paymentInfo.cvv,
            },
        },
    };
    switch (task?.mode.toUpperCase()) {
        case "NORMAL":
            if (task.store.longCode === 'FOOTLOCKERCA') {
                const myFootSite = new FootsitesCA.Footsites(task1, mainWindow)
                myFootSite.start()
            } else {
                const myFootSite = new Footsites.Footsites(task1, mainWindow)
                myFootSite.start()
            }
            break;
        case "PROXY-KILLER":
            if (task.store.longCode === 'FOOTLOCKERCA') {
                const myFootSite = new FootsitesProxyKillerCA.Footsites(task1, mainWindow)
                   myFootSite.start()
            } else {
                const myFootSite = new FootsitesProxyKiller.Footsites(task1, mainWindow)
                myFootSite.start()
            }
            break;
        case "ADVANCED":
            if (task.store.longCode === 'FOOTLOCKERCA') {
                const myFootSite = new FootsitesProxyAdvancedCA.Footsites(task1, mainWindow)
                   myFootSite.start()
            } else {
                const myFootSite = new FootsitesProxyAdvanced.Footsites(task1, mainWindow)
                myFootSite.start()
            }
            break;
        case "WISE":
            if (task.store.longCode === 'FOOTLOCKERCA') {
                
            } else {
                const myFootSite = new FootsitesWise.Footsites(task1, mainWindow)
                myFootSite.start()
            }
            break;
        default:
            break;
    }
}

module.exports = runFootsites;


function abbrState(input, to) {
    var states = [
        ['Arizona', 'AZ'],
        ['Alabama', 'AL'],
        ['Alaska', 'AK'],
        ['Arkansas', 'AR'],
        ['California', 'CA'],
        ['Colorado', 'CO'],
        ['Connecticut', 'CT'],
        ['Delaware', 'DE'],
        ['Florida', 'FL'],
        ['Georgia', 'GA'],
        ['Hawaii', 'HI'],
        ['Idaho', 'ID'],
        ['Illinois', 'IL'],
        ['Indiana', 'IN'],
        ['Iowa', 'IA'],
        ['Kansas', 'KS'],
        ['Kentucky', 'KY'],
        ['Louisiana', 'LA'],
        ['Maine', 'ME'],
        ['Maryland', 'MD'],
        ['Massachusetts', 'MA'],
        ['Michigan', 'MI'],
        ['Minnesota', 'MN'],
        ['Mississippi', 'MS'],
        ['Missouri', 'MO'],
        ['Montana', 'MT'],
        ['Nebraska', 'NE'],
        ['Nevada', 'NV'],
        ['New Hampshire', 'NH'],
        ['New Jersey', 'NJ'],
        ['New Mexico', 'NM'],
        ['New York', 'NY'],
        ['North Carolina', 'NC'],
        ['North Dakota', 'ND'],
        ['Ohio', 'OH'],
        ['Oklahoma', 'OK'],
        ['Oregon', 'OR'],
        ['Pennsylvania', 'PA'],
        ['Rhode Island', 'RI'],
        ['South Carolina', 'SC'],
        ['South Dakota', 'SD'],
        ['Tennessee', 'TN'],
        ['Texas', 'TX'],
        ['Utah', 'UT'],
        ['Vermont', 'VT'],
        ['Virginia', 'VA'],
        ['Washington', 'WA'],
        ['West Virginia', 'WV'],
        ['Wisconsin', 'WI'],
        ['Wyoming', 'WY'],
        ['Alberta', 'AB'],
        ['British Columbia', 'BC'],
        ['Manitoba', 'MB'],
        ['New Brunswick', 'NB'],
        ['Newfoundland and Labrador', 'NL'],
        ['Northwest Territories', 'NT'],
        ['Nova Scotia', 'NS'],
        ['Nunavut', 'NU'],
        ['Ontario', 'ON'],
        ['Prince Edward Island', 'PE'],
        ['Quebec', 'QC'],
        ['Saskatchewan', 'SK'],
        ['Yukon Territory', 'YT']
    ];

    if (to == 'abbr') {
        input = input.replace(/\w\S*/g, function (txt) {
            return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
        });
        for (var i = 0; i < states.length; i++) {
            if (states[i][0] == input) {
                return (states[i][1]);
            }
        }
    } else if (to == 'name') {
        input = input.toUpperCase();
        for (i = 0; i < states.length; i++) {
            if (states[i][1] == input) {
                return (states[i][0]);
            }
        }
    }
}