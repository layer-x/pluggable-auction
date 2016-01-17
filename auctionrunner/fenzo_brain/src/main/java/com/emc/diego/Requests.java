package com.emc.diego;

public class Requests {

    public static class Resource {
        public int MemoryMB;
        public int DiskMB;
    }

    public static class Resources {
        public int MemoryMB;
        public int DiskMB;
        public int Containers;
    }

    public static class LRP {
        public String ProcessGuid;
        public int Index;
        public String Domain;
        public String[] Tags;
        public int MemoryMB;
        public int DiskMB;
        @Override
        public String toString(){
            return "Guid: " + ProcessGuid
                    + "\n" + "DiskMB: " + DiskMB
                    + "\n" + "MemoryMB: " + MemoryMB;
        }
    }

    public static class Task {
        public String TaskGuid;
        public String Domain;
        public String[] Tags;
        public Resource Resource;
        @Override
        public String toString() {
            return "Guid: " + TaskGuid
                    + "\n" + "DiskMB: " + Resource.DiskMB
                    + "\n" + "MemoryMB: " + Resource.MemoryMB;
        }
    }

    public static class SerializableCellState {
        public String id;
        public Resources AvailableResources;
        public Resources TotalResources;
        public LRP[] LRPs;
        public Task[] Tasks;
        public String Zone;
        public Boolean Evacuating;
        public String Guid;
        @Override
        public String toString() {
            return "id: " + id
                    + "\n" + "Guid: " + Guid
                    + "\n" + "DiskMB: " + AvailableResources.DiskMB
                    + "\n" + "MemoryMB: " + AvailableResources.MemoryMB
                    + "\n" + "Containers: " + AvailableResources.Containers;
        }
    }

    public static class AuctionLRPRequest {
        public SerializableCellState[] SerializableCellStates;
        public LRP LRP;
    }

    public static class AuctionTaskRequest {
        public SerializableCellState[] SerializableCellStates;
        public Task Task;
    }
}
