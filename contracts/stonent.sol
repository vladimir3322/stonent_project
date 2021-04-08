pragma solidity ^0.6.0;

import "https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/evm-contracts/src/v0.6/ChainlinkClient.sol";

interface IERC20 {
    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev Moves `amount` tokens from the caller's account to `recipient`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transfer(address recipient, uint256 amount) external returns (bool);

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve} or {transferFrom} are called.
     */
    function allowance(address owner, address spender) external view returns (uint256);

    /**
     * @dev Sets `amount` as the allowance of `spender` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * IMPORTANT: Beware that changing an allowance with this method brings the risk
     * that someone may use both the old and the new allowance by unfortunate
     * transaction ordering. One possible solution to mitigate this race
     * condition is to first reduce the spender's allowance to 0 and set the
     * desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     *
     * Emits an {Approval} event.
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev Moves `amount` tokens from `sender` to `recipient` using the
     * allowance mechanism. `amount` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    /**
     * @dev Emitted when `value` tokens are moved from one account (`from`) to
     * another (`to`).
     *
     * Note that `value` may be zero.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev Emitted when the allowance of a `spender` for an `owner` is set by
     * a call to {approve}. `value` is the new allowance.
     */
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

abstract contract Context {
    function _msgSender() internal view virtual returns (address payable) {
        return msg.sender;
    }

    function _msgData() internal view virtual returns (bytes memory) {
        this; // silence state mutability warning without generating bytecode - see https://github.com/ethereum/solidity/issues/2691
        return msg.data;
    }
}

/**
 * @dev Contract module which provides a basic access control mechanism, where
 * there is an account (an owner) that can be granted exclusive access to
 * specific functions.
 *
 * By default, the owner account will be the one that deploys the contract. This
 * can later be changed with {transferOwnership}.
 *
 * This module is used through inheritance. It will make available the modifier
 * `onlyOwner`, which can be applied to your functions to restrict their use to
 * the owner.
 */
abstract contract Ownable is Context {
    address private _owner;

    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    /**
     * @dev Initializes the contract setting the deployer as the initial owner.
     */
    constructor () public {
        address msgSender = _msgSender();
        _owner = msgSender;
        emit OwnershipTransferred(address(0), msgSender);
    }

    /**
     * @dev Returns the address of the current owner.
     */
    function owner() public view virtual returns (address) {
        return _owner;
    }

    /**
     * @dev Throws if called by any account other than the owner.
     */
    modifier onlyOwner() {
        require(owner() == _msgSender(), "Ownable: caller is not the owner");
        _;
    }

    /**
     * @dev Leaves the contract without owner. It will not be possible to call
     * `onlyOwner` functions anymore. Can only be called by the current owner.
     *
     * NOTE: Renouncing ownership will leave the contract without an owner,
     * thereby removing any functionality that is only available to the owner.
     */
    function renounceOwnership() public virtual onlyOwner {
        emit OwnershipTransferred(_owner, address(0));
        _owner = address(0);
    }

    /**
     * @dev Transfers ownership of the contract to a new account (`newOwner`).
     * Can only be called by the current owner.
     */
    function transferOwnership(address newOwner) public virtual onlyOwner {
        require(newOwner != address(0), "Ownable: new owner is the zero address");
        emit OwnershipTransferred(_owner, newOwner);
        _owner = newOwner;
    }
}

contract Stonent is ChainlinkClient, Ownable {
    address public oracle;
    bytes32 public jobId;
    uint256 public fee;
    uint256 public version;
    uint256 public price;
    address private dev;
    address public paymentToken;

    mapping (bytes32 => Certificate) public certificates;
    mapping (string => bytes32) public lastCertification;

    struct Certificate {
        uint256 Score;
        address Oracle;
        uint256 Version;
        uint256 Date;
        string ID;
    }

    event RequestSended(address indexed add,string indexed id, address token, bytes32 reqID);
    event RequestCertified(uint256 indexed score, bytes32 reqID);

    uint256 public result;


    constructor() public {
        setPublicChainlinkToken();
        oracle = 0xc00AF972A21f783931A4d558C0996Ae66CE7E505;
        jobId = "55ab785f7b424544afc45390690649c8";
        fee = 1 * 10 ** 18;
        dev = msg.sender;
        price = 100 * 10 ** 18; // price for certification in paymentToken

    }

    /**
     * Initial request
     */
    function check(string memory _id) public {

        IERC20(paymentToken).transferFrom(msg.sender, dev, price);

        Chainlink.Request memory req = buildChainlinkRequest(jobId, address(this), this.fulfillScore.selector);
        req.add("id", _id);
        bytes32 requestID = sendChainlinkRequestTo(oracle, req, fee);

        Certificate storage certificate = certificates[requestID];

        certificate.ID = _id;
        certificate.Version = version;
        certificate.Oracle = oracle;

        lastCertification[_id] = requestID;


        emit RequestSended(msg.sender, _id, address(0x0), requestID);
    }

    /**
     * Callback function
     */
    function fulfillScore(bytes32 _requestId, uint256 _score) public recordChainlinkFulfillment(_requestId) {
        Certificate storage certificate = certificates[_requestId];

        certificate.Score = _score;
        certificate.Date = now;

         emit RequestCertified(_score, _requestId);
    }

    function changeOracle(address _oracle) onlyOwner public {
        oracle = _oracle;
    }

    function changeJobID(bytes32 _job) onlyOwner public {
        jobId = _job;
    }

    function changeOracleFee(uint256 _fee) onlyOwner public {
        fee = _fee;
    }

    function changePrice(uint256 _price) onlyOwner public {
        price = _price;
    }

    function changeDev(address _newDev) onlyOwner public {
        dev = _newDev;
    }

    function changeAdapterVersion(uint256 _v) onlyOwner public {
        version = _v;
    }

    function setPaymentToken(address _token) onlyOwner public {
        paymentToken = _token;
    }
}
