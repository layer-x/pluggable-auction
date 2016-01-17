package com.emc.diego;

import com.google.gson.Gson;
import com.netflix.fenzo.*;
import com.netflix.fenzo.functions.Action1;
import spark.Spark;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class FenzoBrain
{
    public static TaskScheduler myTaskScheduler;
    public static Map<String, Requests.SerializableCellState> knownCells;
    public static void main( String[] args )
    {
        knownCells = new HashMap<>();
        System.out.println("Starting TaskScheduler...");
        myTaskScheduler = new TaskScheduler.Builder()
                .withLeaseRejectAction(new Action1<VirtualMachineLease>() {
                    @Override
                    public void call(VirtualMachineLease virtualMachineLease) {
                        System.out.println("Declining offer on " + virtualMachineLease.hostname());
                    }
                })
                .withLeaseOfferExpirySecs(1000000000)
                .build();

        System.out.println( "Starting FenzoBrain on port 5555" );

        Spark.port(5555);
        Spark.post("/AuctionLRP", (req, res) -> {
            String auctionLRPRequestJson = req.body();
            Gson gson = new Gson();
            Requests.AuctionLRPRequest auctionLRPRequest = gson.fromJson(auctionLRPRequestJson, Requests.AuctionLRPRequest.class);
            byte[] resData = auctionLRP(auctionLRPRequest);
            if (resData.length == 0) {
                res.status(500);
            } else {
                res.body(new String(resData));
            }
            return new String(resData);
        });
        Spark.post("/AuctionTask", (req, res) -> {
            String auctionTaskRequestJson = req.body();
            Gson gson = new Gson();
            Requests.AuctionTaskRequest auctionTaskRequest = gson.fromJson(auctionTaskRequestJson, Requests.AuctionTaskRequest.class);
            byte[] resData = auctionTask(auctionTaskRequest);
            if (resData.length == 0) {
                res.status(500);
            } else {
                res.body(new String(resData));
            }
            return new String(resData);
        });
    }

    public static byte[] auctionLRP(Requests.AuctionLRPRequest auctionLRPRequest) {
        Requests.SerializableCellState winnerCell = null;
        List<VirtualMachineLease> newLeases = new ArrayList<>();
        System.out.println("Trying to schedule LRP: "+auctionLRPRequest.LRP.toString()+" on cells: ");
        for (Requests.SerializableCellState cellState : auctionLRPRequest.SerializableCellStates) {
            System.out.println(cellState.toString());
            if (!knownCells.containsKey(cellState.id)) {
                VirtualMachineLease newLease = VmLeaseFactory.fromCellState(cellState);
                newLeases.add(newLease);
            }
            knownCells.put(cellState.id, cellState);
        }
        List<TaskRequest> lrpRequests = new ArrayList<>();
        lrpRequests.add(TaskRequestFactory.fromLRP(auctionLRPRequest.LRP));
        SchedulingResult result = null;
        try {
            result = myTaskScheduler.scheduleOnce(lrpRequests, newLeases);
        } catch (IllegalStateException e) {
            System.out.println("Expiring all leases...");
            knownCells.clear();
            myTaskScheduler.expireAllLeases();
        }
        System.out.println("result=" + result);
        if (result == null) {
            return new byte[0];
        }
        Map<String, VMAssignmentResult> resultMap = result.getResultMap();
        System.out.println("about to do stuff with result...");
        if (!resultMap.isEmpty()) {
            System.out.println("result is not empty");
            for (VMAssignmentResult vmAssignmentResult : resultMap.values()) {
                System.out.println("parsing vmAssignment Result..." + vmAssignmentResult);
                List<VirtualMachineLease> leasesUsed = vmAssignmentResult.getLeasesUsed();
                for (VirtualMachineLease lease : leasesUsed) {
                    for (TaskAssignmentResult taskAssignmentResult : vmAssignmentResult.getTasksAssigned()) {
                        if (taskAssignmentResult.getRequest() != null) {
                            System.out.println("We found a match! Thank god: LRP " + taskAssignmentResult.getRequest().getId() + " assigned to Cell " + lease.getVMID());
                            myTaskScheduler.getTaskAssigner().call(taskAssignmentResult.getRequest(), lease.hostname());
                            String cellId = lease.getId();
                            winnerCell = knownCells.get(cellId);
                        }
                    }
                }
            }
        }
        if (winnerCell == null) {
            System.out.println("DID NOT ASSIGN THE LRP! WHY!!!");
            return new byte[0];
        }
        Gson gson = new Gson();
        String winnerJson = gson.toJson(winnerCell);
        return winnerJson.getBytes();
    }

    public static byte[] auctionTask(Requests.AuctionTaskRequest auctionTaskRequest) {
        Requests.SerializableCellState winnerCell = null;
        List<VirtualMachineLease> newLeases = new ArrayList<>();
        for (Requests.SerializableCellState cellState : auctionTaskRequest.SerializableCellStates) {
            if (!knownCells.containsKey(cellState.id)) {
                VirtualMachineLease newLease = VmLeaseFactory.fromCellState(cellState);
                newLeases.add(newLease);
            }
            knownCells.put(cellState.id, cellState);
        }
        List<TaskRequest> taskRequests = new ArrayList<>();
        taskRequests.add(TaskRequestFactory.fromTask(auctionTaskRequest.Task));
        SchedulingResult result = null;
        try {
            result = myTaskScheduler.scheduleOnce(taskRequests, newLeases);
        } catch (IllegalStateException e) {
            System.out.println("Expiring all leases...");
            knownCells.clear();
            myTaskScheduler.expireAllLeases();
        }
        System.out.println("result=" + result);
        if (result == null) {
            return new byte[0];
        }
        Map<String, VMAssignmentResult> resultMap = result.getResultMap();
        System.out.println("about to do stuff with result...");
        if (!resultMap.isEmpty()) {
            System.out.println("result is not empty");
            for (VMAssignmentResult vmAssignmentResult : resultMap.values()) {
                System.out.println("parsing vmAssignment Result..." + vmAssignmentResult);
                List<VirtualMachineLease> leasesUsed = vmAssignmentResult.getLeasesUsed();
                for (VirtualMachineLease lease : leasesUsed) {
                    for (TaskAssignmentResult taskAssignmentResult : vmAssignmentResult.getTasksAssigned()) {
                        if (taskAssignmentResult.getRequest() != null) {
                            System.out.println("We found a match! Thank god: Task " + taskAssignmentResult.getRequest().getId() + " assigned to Cell " + lease.getVMID());
                            myTaskScheduler.getTaskAssigner().call(taskAssignmentResult.getRequest(), lease.hostname());
                            String cellId = lease.getId();
                            winnerCell = knownCells.get(cellId);
                        }
                    }
                }
            }
        }
        if (winnerCell == null) {
            System.out.println("DID NOT ASSIGN THE Task! WHY!!!");
            return new byte[0];
        }
        Gson gson = new Gson();
        String winnerJson = gson.toJson(winnerCell);
        return winnerJson.getBytes();
    }
}
