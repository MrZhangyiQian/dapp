pragma solidity ^0.8.7;

contract Voting {
    // mapping 来存储 候选人和其对应得票数
    mapping(address => uint256) public votes;
    
    function vote(address candidate, bool support) public {
        if (support == true){
            votes[candidate]++;
        } else{
           require(votes[candidate] > 0,"You cant remove a vote that wasnt given");
           votes[candidate]--;
       }
   }

function getVotes(address candidate) public view returns(uint256){
    return votes[candidate];
}

// 重置所有候选人的得票数
function resetVoteForAddress(address addressToReset) public {
      require(votes[addressToReset] > 0,"You cant remove a vote that wasnt given");
   votes[addressToReset] = 0;
}
}