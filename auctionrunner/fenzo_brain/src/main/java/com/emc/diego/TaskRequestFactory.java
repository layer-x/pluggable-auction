package com.emc.diego;

import com.netflix.fenzo.ConstraintEvaluator;
import com.netflix.fenzo.TaskRequest;
import com.netflix.fenzo.VMTaskFitnessCalculator;

import java.util.List;
import java.util.UUID;

/**
 * Created by pivotal on 1/16/16.
 */
public class TaskRequestFactory {
    public static TaskRequest fromLRP(Requests.LRP lrp) {
        return new TaskRequest() {
            @Override
            public String getId() {
                UUID id = UUID.randomUUID();
                return id.toString();
            }

            @Override
            public String taskGroupName() {
                return null;
            }

            @Override
            public double getCPUs() {
                return 1;
            }

            @Override
            public double getMemory() {
                return lrp.MemoryMB;
            }

            @Override
            public double getNetworkMbps() {
                return 0;
            }

            @Override
            public double getDisk() {
                return lrp.DiskMB;
            }

            @Override
            public int getPorts() {
                return 0;
            }

            @Override
            public List<? extends ConstraintEvaluator> getHardConstraints() {
                return null;
            }

            @Override
            public List<? extends VMTaskFitnessCalculator> getSoftConstraints() {
                return null;
            }
        };
    }
    public static TaskRequest fromTask(Requests.Task task) {
        return new TaskRequest() {
            @Override
            public String getId() {
                UUID id = UUID.randomUUID();
                return id.toString();
            }

            @Override
            public String taskGroupName() {
                return null;
            }

            @Override
            public double getCPUs() {
                return 0;
            }

            @Override
            public double getMemory() {
                return task.Resource.MemoryMB;
            }

            @Override
            public double getNetworkMbps() {
                return 0;
            }

            @Override
            public double getDisk() {
                return task.Resource.DiskMB;
            }

            @Override
            public int getPorts() {
                return 0;
            }

            @Override
            public List<? extends ConstraintEvaluator> getHardConstraints() {
                return null;
            }

            @Override
            public List<? extends VMTaskFitnessCalculator> getSoftConstraints() {
                return null;
            }
        };
    }
}
