package com.emc.diego;

import com.netflix.fenzo.VirtualMachineLease;
import org.apache.mesos.Protos;

import java.util.List;
import java.util.Map;

/**
 * Created by pivotal on 1/16/16.
 */
public class VmLeaseFactory {
    public static VirtualMachineLease fromCellState(Requests.SerializableCellState cell) {
        return new VirtualMachineLease() {
            @Override
            public String getId() {
                return cell.id;
            }

            @Override
            public long getOfferedTime() {
                return 0;
            }

            @Override
            public String hostname() {
                return cell.Guid;
            }

            @Override
            public String getVMID() {
                return cell.Guid;
            }

            @Override
            public double cpuCores() {
                return cell.AvailableResources.Containers;
            }

            @Override
            public double memoryMB() {
                return cell.AvailableResources.MemoryMB;
            }

            @Override
            public double networkMbps() {
                return 0;
            }

            @Override
            public double diskMB() {
                return cell.AvailableResources.DiskMB;
            }

            @Override
            public List<Range> portRanges() {
                return null;
            }

            @Override
            public Protos.Offer getOffer() {
                return null;
            }

            @Override
            public Map<String, Protos.Attribute> getAttributeMap() {
                return null;
            }
        };
    }
}
