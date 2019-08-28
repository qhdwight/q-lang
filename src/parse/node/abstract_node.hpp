#pragma once

#include <memory>
#include <vector>
#include <utility>

namespace ql::parse {
    class AbstractNode {
    public:
        typedef std::vector<std::shared_ptr<AbstractNode>> ChildrenRef;
        typedef std::weak_ptr<AbstractNode> ParentRef;
    private:
        ChildrenRef m_Children;
        ParentRef m_Parent;
    public:
        explicit AbstractNode(ParentRef parent) : m_Parent(std::move(parent)) {}

        void addChild(const std::shared_ptr<AbstractNode>& node);
    };
}
