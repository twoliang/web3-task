// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract Voting {
    mapping(string => uint256) public voteCount;

    string[] public candidates;

    event Voted(string candidate, uint256 newVoteCount);

    event VotesReset(string candidate);
    event AllVotesReset();

    function vote(string memory _candidate) public {
        if (voteCount[_candidate] == 0) {
            candidates.push(_candidate);
        }
        voteCount[_candidate]++;
        emit Voted(_candidate, voteCount[_candidate]);
    }

    function getVotes(string memory _candidate) public view returns (uint256) {
        return voteCount[_candidate];
    }

    function resetVotes(string memory _candidate) public {
        if (bytes(_candidate).length == 0) {
            for (uint i = 0; i < candidates.length; i++) {
                voteCount[candidates[i]] = 0;
            }
            emit AllVotesReset();
        } 
        else {
            voteCount[_candidate] = 0;
            emit VotesReset(_candidate);
        }
    }
}