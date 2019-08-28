#pragma once

#include "abstract_node.hpp"

namespace ql::parse {
    class MasterNode : public AbstractNode {
    public:
        MasterNode() : AbstractNode(ParentRef()) {}
    };
}